package generator

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/builder"
	"github.com/jmattheis/goverter/config"
	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/namer"
	"github.com/jmattheis/goverter/xtype"
)

type generatedMethod struct {
	*config.Method

	Explicit bool
	Dirty    bool

	OriginPath []method.IndexID
	Jen        jen.Code

	IndexID method.IndexID
}

type generator struct {
	namer  *namer.Namer
	conf   *config.Converter
	lookup *method.Index[generatedMethod]
	extend *method.Index[method.Definition]
}

func (g *generator) getGenMethods() []*generatedMethod {
	genMethods := g.lookup.GetAll()
	sort.Slice(genMethods, func(i, j int) bool {
		return genMethods[i].Name < genMethods[j].Name
	})
	return genMethods
}

func (g *generator) buildMethods(f *jen.File) error {
	for g.anyDirty() {
		if err := g.buildDirtyMethods(); err != nil {
			return err
		}
	}
	g.appendGenerated(f)
	return nil
}

func (g *generator) buildDirtyMethods() error {
	for _, genMethod := range g.getGenMethods() {
		if !genMethod.Dirty {
			continue
		}
		genMethod.Dirty = false
		err := g.buildMethod(genMethod, genMethod.Context)
		if err != nil {
			err = err.Lift(&builder.Path{
				SourceID:   "source",
				TargetID:   "target",
				SourceType: genMethod.Source.String,
				TargetType: genMethod.Target.String,
			})
			return fmt.Errorf("Error while creating converter method:\n    %s\n    %s%s\n\n%s", genMethod.Location, genMethod.ID, genMethod.Definition.ArgDebug("        "), builder.ToString(err))
		}
	}
	return nil
}

func (g *generator) anyDirty() bool {
	for _, m := range g.getGenMethods() {
		if m.Dirty {
			return true
		}
	}
	return false
}

func (g *generator) appendGenerated(f *jen.File) {
	genMethods := g.getGenMethods()
	for _, raw := range g.conf.OutputRaw {
		f.Id(raw)
	}

	if g.conf.OutputFormat == config.FormatStruct {
		if len(g.conf.Comments) > 0 {
			f.Comment(strings.Join(g.conf.Comments, "\n"))
		}
		f.Type().Id(g.conf.Name).Struct()
	}

	var init []jen.Code
	var funcs []jen.Code

	for _, def := range genMethods {
		switch g.conf.OutputFormat {
		case config.FormatStruct:
			funcs = append(funcs, jen.Func().Params(jen.Id(xtype.ThisVar).Op("*").Id(g.conf.Name)).Id(def.Name).Add(def.Jen))
		case config.FormatVariable:
			if def.Explicit {
				init = append(init, jen.Qual(def.Package, def.Name).Op("=").Func().Add(def.Jen))
			} else {
				funcs = append(funcs, jen.Func().Id(def.Name).Add(def.Jen))
			}
		case config.FormatFunction:
			funcs = append(funcs, jen.Func().Id(def.Name).Add(def.Jen))
		}
	}

	if len(init) > 0 {
		f.Func().Id("init").Params().Block(init...)
	}

	for _, fn := range funcs {
		f.Add(fn)
	}
}

func (g *generator) buildMethod(genMethod *generatedMethod, context map[string]*xtype.Type) *builder.Error {
	var sourceID *xtype.JenID
	source := genMethod.Source
	target := genMethod.Target

	fieldsTarget := genMethod.Target.String
	if genMethod.Target.Pointer && genMethod.Target.PointerInner.Struct {
		fieldsTarget = genMethod.Target.PointerInner.String
	}

	ctx := &builder.MethodContext{
		Namer:             namer.New(),
		Conf:              genMethod.Method,
		FieldsTarget:      fieldsTarget,
		AvailableContext:  context,
		SeenNamed:         map[string]struct{}{},
		TargetType:        genMethod.Target,
		Context:           map[string]*xtype.JenID{},
		IndexID:           genMethod.IndexID,
		Signature:         genMethod.Signature,
		HasMethod:         g.hasMethod,
		OutputPackagePath: g.conf.OutputPackagePath,
		UseConstructor:    genMethod.Constructor != nil,
	}

	var targetAssign *jen.Statement
	args := []jen.Code{}
	for _, arg := range genMethod.RawArgs {
		switch arg.Use {
		case method.ArgUseInterface:
			panic("hopefully unreachable")
		case method.ArgUseContext:
			name := ctx.Name("context")
			ctx.Context[arg.Type.String] = xtype.VariableID(jen.Id(name))
			args = append(args, jen.Id(name).Add(arg.Type.TypeAsJen()))
		case method.ArgUseSource:
			name := ctx.Name("source")
			sourceID = xtype.VariableID(jen.Id(name))
			args = append(args, jen.Id(name).Add(arg.Type.TypeAsJen()))
		case method.ArgUseTarget:
			name := ctx.Name("target")
			targetAssign = jen.Id(name)
			args = append(args, jen.Id(name).Add(arg.Type.TypeAsJen()))
		case method.ArgUseMultiSource:
			panic("multi source aren't supported right now. https://github.com/jmattheis/goverter/issues/143")
		}
	}

	var returns []jen.Code
	if targetAssign == nil {
		returns = append(returns, target.TypeAsJen())
	}

	if genMethod.ReturnError {
		returns = append(returns, jen.Id("error"))
	}

	var funcBlock []jen.Code
	if targetAssign != nil {
		var err *builder.Error
		funcBlock, err = g.convertTo(ctx, builder.AssignOf(targetAssign), sourceID, source, target, nil)
		if err != nil {
			return err
		}

		if genMethod.ReturnError {
			funcBlock = append(funcBlock, jen.Return().Nil())
		}
	} else if def, err := g.extend.Get(ctx.Signature, context); def != nil {
		jenReturn, err := g.delegateMethod(ctx, def, sourceID)
		if err != nil {
			return err
		}
		funcBlock = []jen.Code{jenReturn}
	} else if err != nil {
		return builder.NewError(err.Error())
	} else {
		stmt, newID, err := g.buildNoLookup(ctx, sourceID, source, target, nil)
		if err != nil {
			return err
		}
		ret := []jen.Code{newID.Code}
		if genMethod.ReturnError {
			ret = append(ret, jen.Nil())
		}

		funcBlock = append(stmt, jen.Return(ret...))
	}

	genMethod.Jen = jen.Params(args...).Params(returns...).Block(funcBlock...)

	return nil
}

func (g *generator) buildNoLookup(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPath builder.ErrorPath) ([]jen.Code, *xtype.JenID, *builder.Error) {
	if err := g.getOverlappingStructDefinition(ctx, source, target); err != nil {
		return nil, nil, err
	}

	for _, rule := range BuildSteps {
		if rule.Matches(ctx, source, target) {
			return rule.Build(g, ctx, sourceID, source, target, errPath)
		}
	}

	return nil, nil, typeMismatch(source, target)
}

func (g *generator) assignNoLookup(ctx *builder.MethodContext, assignTo *builder.AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath builder.ErrorPath) ([]jen.Code, *builder.Error) {
	if err := g.getOverlappingStructDefinition(ctx, source, target); err != nil {
		return nil, err
	}

	for _, rule := range BuildSteps {
		if rule.Matches(ctx, source, target) {
			return rule.Assign(g, ctx, assignTo, sourceID, source, target, errPath)
		}
	}

	return nil, typeMismatch(source, target)
}

func (g *generator) convertTo(ctx *builder.MethodContext, assignTo *builder.AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath builder.ErrorPath) ([]jen.Code, *builder.Error) {
	if !target.Pointer || !target.PointerInner.Struct {
		return nil, builder.NewError("target type must be a pointer struct for goverter:update signatures.")
	}
	sourcePointer := false
	if !source.Struct {
		if source.Pointer && source.PointerInner.Struct {
			sourcePointer = true
			source = source.PointerInner
		} else {
			return nil, builder.NewError("source type must be a struct or pointer struct for goverter:update signatures.")
		}
	}

	var s builder.Struct
	stmt, err := s.Assign(g, ctx, assignTo, sourceID, source, target.PointerInner, errPath)
	if sourcePointer {
		stmt = []jen.Code{jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(stmt...)}
	}
	return stmt, err
}

func (g *generator) CallMethod(
	ctx *builder.MethodContext,
	definition *method.Definition,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errPath builder.ErrorPath,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	params := []jen.Code{}
	formatErr := func(s string) *builder.Error {
		return builder.NewError(fmt.Sprintf("Error using method:\n    %s%s\n\n%s", definition.ID, definition.ArgDebug("        "), s))
	}

	for _, arg := range definition.RawArgs {
		switch arg.Use {
		case method.ArgUseInterface:
			params = append(params, jen.Id(xtype.ThisVar))
		case method.ArgUseContext:
			if !g.requireContext(ctx, arg.Type) {
				return nil, nil, formatErr("Could not satisfy all required context parameters:\n" + strings.Join(method.AvailableContextDebug(definition.Context, ctx.AvailableContext), "\n"))
			}
			if id, ok := ctx.Context[arg.Type.String]; ok {
				params = append(params, id.Code.Clone())
			}
		case method.ArgUseSource:
			if !source.AssignableTo(definition.Source) && !definition.TypeParams {
				cause := fmt.Sprintf("Method source type mismatches with conversion source: %s != %s", definition.Source.String, source.String)
				return nil, nil, formatErr(cause)
			}
			params = append(params, sourceID.Code)
		case method.ArgUseMultiSource:
			panic("multi source aren't supported right now. https://github.com/jmattheis/goverter/issues/143")
		case method.ArgUseTarget:
			panic("unreachable")
		}
	}

	if !definition.Target.AssignableTo(target) && !definition.TypeParams {
		cause := fmt.Sprintf("Method return type mismatches with target: %s != %s", definition.Target.String, target.String)
		return nil, nil, formatErr(cause)
	}

	qual := g.qualMethod(definition)
	if definition.ReturnError {
		name := ctx.Name(target.ID())
		ctx.SetErrorTargetVar(jen.Id(name))

		ret, ok := g.ReturnError(ctx, errPath, jen.Id("err"))
		if !ok {
			return nil, nil, formatErr("Used method returns error but conversion method does not")
		}

		stmt := []jen.Code{
			jen.List(jen.Id(name), jen.Id("err")).Op(":=").Add(qual.Call(params...)),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(ret),
		}
		return stmt, xtype.VariableID(jen.Id(name)), nil
	}
	id := xtype.OtherID(qual.Call(params...))
	return nil, id, nil
}

func (g *generator) ReturnError(ctx *builder.MethodContext, errPath builder.ErrorPath, id *jen.Statement) (jen.Code, bool) {
	current := g.lookup.ByID(ctx.IndexID)
	if !ctx.Conf.ReturnError {
		for _, path := range append([]method.IndexID{ctx.IndexID}, current.OriginPath...) {
			check := g.lookup.ByID(path)
			if check.Explicit && !check.ReturnError {
				return nil, false
			}

			if !check.ReturnError {
				check.ReturnError = true
				check.Dirty = true
			}
		}
	}
	returns := []jen.Code{}
	if !current.UpdateTarget {
		returns = append(returns, ctx.TargetVar)
	}
	returns = append(returns, g.wrap(ctx, errPath, id))
	return jen.Return(returns...), true
}

func (g *generator) requireContext(ctx *builder.MethodContext, need *xtype.Type) bool {
	if _, ok := ctx.Context[need.String]; ok {
		return true
	}

	current := g.lookup.ByID(ctx.IndexID)
	for _, path := range append([]method.IndexID{ctx.IndexID}, current.OriginPath...) {
		check := g.lookup.ByID(path)

		if _, ok := check.Context[need.String]; ok {
			continue
		}

		if check.Explicit {
			return false
		}

		check.Context[need.String] = need
		check.RawArgs = append(check.RawArgs, method.Arg{
			Name: "",
			Use:  method.ArgUseContext,
			Type: need,
		})
		check.Dirty = true
	}
	return true
}

func (g *generator) delegateMethod(
	ctx *builder.MethodContext,
	delegateTo *method.Definition,
	sourceID *xtype.JenID,
) (*jen.Statement, *builder.Error) {
	params := []jen.Code{}

	for _, arg := range delegateTo.RawArgs {
		switch arg.Use {
		case method.ArgUseInterface:
			params = append(params, jen.Id(xtype.ThisVar))
		case method.ArgUseContext:
			params = append(params, ctx.Context[arg.Type.String].Code.Clone())
		case method.ArgUseSource:
			params = append(params, sourceID.Code)
		case method.ArgUseMultiSource:
			panic("not supported atm")
		case method.ArgUseTarget:
			panic("unreachable")
		}
	}

	current := g.lookup.ByID(ctx.IndexID)

	returns := []jen.Code{g.qualMethod(delegateTo).Call(params...)}

	if delegateTo.ReturnError {
		if !current.ReturnError {
			return nil, builder.NewError(fmt.Sprintf("ReturnTypeMismatch: Cannot use\n\n    %s\n\nin\n\n    %s\n\nbecause no error is returned as second return parameter", delegateTo.OriginID, current.ID))
		}
	} else {
		if current.ReturnError {
			returns = append(returns, jen.Nil())
		}
	}
	return jen.Return(returns...), nil
}

// wrap invokes the error wrapper if feature is enabled.
func (g *generator) wrap(ctx *builder.MethodContext, errPath builder.ErrorPath, errStmt *jen.Statement) *jen.Statement {
	switch {
	case ctx.Conf.WrapErrorsUsing != "":
		return errPath.WrapErrorsUsing(ctx.Conf.WrapErrorsUsing, errStmt)
	case ctx.Conf.WrapErrors:
		return errPath.WrapErrors(errStmt)
	default:
		return errStmt
	}
}

// Build builds an implementation for the given source and target type, or uses an existing method for it.
func (g *generator) Build(
	ctx *builder.MethodContext,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errPath builder.ErrorPath,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	stmt, nextID, err := g.callExisting(ctx, sourceID, source, target, errPath)
	if nextID != nil || err != nil {
		return stmt, nextID, err
	}

	if g.shouldCreateSubMethod(ctx, source, target) {
		return g.createSubMethod(ctx, sourceID, source, target, errPath)
	}

	return g.buildNoLookup(ctx, sourceID, source, target, errPath)
}

// Assign builds an implementation for the given source and target type, or uses an existing method for it.
func (g *generator) Assign(
	ctx *builder.MethodContext,
	assignTo *builder.AssignTo,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errPath builder.ErrorPath,
) ([]jen.Code, *builder.Error) {
	if assignTo.Must {
		return builder.ToAssignable(assignTo)(g.Build(ctx, sourceID, source, target, errPath))
	}

	stmt, nextID, err := g.callExisting(ctx, sourceID, source, target, errPath)
	if nextID != nil || err != nil {
		return builder.ToAssignable(assignTo)(stmt, nextID, err)
	}

	if g.shouldCreateSubMethod(ctx, source, target) {
		return builder.ToAssignable(assignTo)(g.createSubMethod(ctx, sourceID, source, target, errPath))
	}

	return g.assignNoLookup(ctx, assignTo, sourceID, source, target, errPath)
}

func (g generator) callExisting(
	ctx *builder.MethodContext,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errPath builder.ErrorPath,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	signature := xtype.SignatureOf(source, target)
	if def, err := g.extend.Get(signature, ctx.AvailableContext); def != nil {
		return g.CallMethod(ctx, def, sourceID, source, target, errPath)
	} else if err != nil {
		return nil, nil, builder.NewError(err.Error())
	}
	if genMethod, err := g.lookup.Get(signature, ctx.AvailableContext); genMethod != nil {
		return g.CallMethod(ctx, genMethod.Definition, sourceID, source, target, errPath)
	} else if err != nil {
		return nil, nil, builder.NewError(err.Error())
	}
	return nil, nil, nil
}

func (g *generator) shouldCreateSubMethod(ctx *builder.MethodContext, source, target *xtype.Type) bool {
	isCurrentPointerStructMethod := false
	if source.Struct && target.Struct {
		// This checks if we are currently inside the generation of one of the following combinations.
		// *Source -> Target
		//  Source -> *Target
		// *Source -> *Target
		isCurrentPointerStructMethod = types.Identical(ctx.Signature.Source, source.AsPointerType()) ||
			types.Identical(ctx.Signature.Target, target.AsPointerType())
	}

	createSubMethod := false

	if ctx.HasSeen(source) {
		g.lookup.ByID(ctx.IndexID).Dirty = true
		createSubMethod = true
	} else if !isCurrentPointerStructMethod {
		switch {
		case source.Named && !source.Basic:
			createSubMethod = true
		case target.Named && !target.Basic:
			createSubMethod = true
		case source.Pointer && source.PointerInner.Named && !source.PointerInner.Basic:
			createSubMethod = true
		case source.Enum(&ctx.Conf.Enum).OK && target.Enum(&ctx.Conf.Enum).OK:
			createSubMethod = true
		}
		if ctx.Conf.SkipCopySameType && types.Identical(source.T, target.T) {
			createSubMethod = false
		}
	}
	ctx.MarkSeen(source)

	return createSubMethod
}

func (g *generator) createSubMethod(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPAth builder.ErrorPath) ([]jen.Code, *xtype.JenID, *builder.Error) {
	name := g.namer.Name(source.UnescapedID() + "To" + strings.Title(target.UnescapedID()))
	orig := g.lookup.ByID(ctx.IndexID)

	var args []method.Arg
	args = append(args, method.Arg{
		Name: "source",
		Type: source,
		Use:  method.ArgUseSource,
	})

	path := append([]method.IndexID{ctx.IndexID}, orig.OriginPath...)
	genMethod := &generatedMethod{
		OriginPath: path,
		Method: &config.Method{
			Common:      g.conf.Common,
			Fields:      map[string]*config.FieldMapping{},
			EnumMapping: &config.EnumMapping{Map: map[string]string{}},
			Definition: &method.Definition{
				OriginID:  ctx.Conf.OriginID,
				ID:        name,
				Package:   g.conf.OutputPackagePath,
				Name:      name,
				Generated: true,
				Parameters: method.Parameters{
					Source:    source,
					RawArgs:   args,
					Context:   map[string]*xtype.Type{},
					Signature: xtype.SignatureOf(source, target),
					Target:    target,
				},
			},
		},
	}

	genMethod.IndexID, _ = g.lookup.Register(genMethod, genMethod.Definition)

	if err := g.buildMethod(genMethod, ctx.AvailableContext); err != nil {
		return nil, nil, err
	}
	return g.CallMethod(ctx, genMethod.Definition, sourceID, source, target, errPAth)
}

func (g *generator) hasMethod(ctx *builder.MethodContext, source, target types.Type) bool {
	signature := xtype.Signature{Source: source, Target: target}
	return g.extend.Has(signature) || g.lookup.Has(signature)
}

func (g *generator) getOverlappingStructDefinition(ctx *builder.MethodContext, source, target *xtype.Type) *builder.Error {
	if !source.Struct || !target.Struct {
		return nil
	}

	overlapping := []xtype.Signature{
		{Source: source.AsPointerType(), Target: target.T},
		{Source: source.AsPointerType(), Target: target.AsPointerType()},
		{Source: source.T, Target: target.AsPointerType()},
	}

	for _, sig := range overlapping {
		if ctx.Signature.Identical(sig) {
			continue
		}
		if def, _ := g.lookup.Get(sig, ctx.AvailableContext); def != nil && len(def.RawFieldSettings) > 0 {
			var toMethod string
			if def, _ := g.lookup.Get(ctx.Signature, ctx.AvailableContext); def != nil && def.Explicit {
				toMethod = fmt.Sprintf("to the %q method.", def.Name)
			} else {
				toMethod = fmt.Sprintf("to a newly created method with this signature:\n    func(%s) %s", source.String, target.String)
			}

			return builder.NewError(fmt.Sprintf(`Overlapping struct settings found.

Move these field related settings:
    goverter:%s

from the %q method %s

Goverter won't use %q inside the current conversion method
and therefore the defined field settings would be ignored.`, strings.Join(def.RawFieldSettings, "\n    goverter:"), def.Name, toMethod, def.Name))
		}
	}
	return nil
}

func typeMismatch(source, target *xtype.Type) *builder.Error {
	if source.Pointer && !target.Pointer {
		return builder.NewError(fmt.Sprintf(`TypeMismatch: Cannot convert %s to %s
It is unclear how nil should be handled in the pointer to non pointer conversion.

You can enable useZeroValueOnPointerInconsistency to instruct goverter to use the zero value if source is nil
https://goverter.jmattheis.de/reference/useZeroValueOnPointerInconsistency

or you can define a custom conversion method with extend:
https://goverter.jmattheis.de/reference/extend`, source.T, target.T))
	}

	return builder.NewError(fmt.Sprintf(`TypeMismatch: Cannot convert %s to %s

You can define a custom conversion method with extend:
https://goverter.jmattheis.de/reference/extend`, source.T, target.T))
}

func (g *generator) qualMethod(m *method.Definition) *jen.Statement {
	switch {
	case m.CustomCall != nil:
		return m.CustomCall.Clone()
	case g.conf.OutputFormat == config.FormatStruct && m.Generated:
		return jen.Id(xtype.ThisVar).Dot(m.Name)
	case g.conf.OutputFormat == config.FormatFunction && m.Generated:
		return jen.Id(m.Name)
	default:
		return jen.Qual(m.Package, m.Name)
	}
}

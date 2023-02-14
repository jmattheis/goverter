package generator

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/builder"
	"github.com/jmattheis/goverter/comments"
	"github.com/jmattheis/goverter/namer"
	"github.com/jmattheis/goverter/xtype"
	"golang.org/x/tools/go/packages"
)

type methodDefinition struct {
	ID       string
	Explicit bool
	Name     string
	Call     *jen.Statement
	Source   *xtype.Type
	Target   *xtype.Type

	Fields          map[string]*builder.FieldMapping
	MatchIgnoreCase bool
	WrapErrors      bool

	Jen jen.Code

	SelfAsFirstParam       bool
	ReturnError            bool
	ReturnTypeOrigin       string
	Dirty                  bool
	IgnoreUnexportedFields bool
}

type generator struct {
	namer  *namer.Namer
	name   string
	file   *jen.File
	lookup map[xtype.Signature]*methodDefinition
	extend map[xtype.Signature]*methodDefinition
	// wrapErrors embeds field names in the error msg
	wrapErrors bool
	// pkgCache caches the extend packages, saving load time
	pkgCache map[string][]*packages.Package
	// workingDir is a working directory, can be empty
	workingDir             string
	ignoreUnexportedFields bool
}

func (g *generator) registerMethod(converter types.Type, scope *types.Scope, methodType *types.Func, methodComments comments.Method) error {
	signature, ok := methodType.Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("expected signature %#v", methodType.Type())
	}
	params := signature.Params()
	if params.Len() != 1 {
		return fmt.Errorf("expected signature to have only one parameter")
	}
	result := signature.Results()
	if result.Len() < 1 || result.Len() > 2 {
		return fmt.Errorf("return has no or too many parameters")
	}
	source := params.At(0).Type()
	target := result.At(0).Type()
	returnError := false
	if result.Len() == 2 {
		if i, ok := result.At(1).Type().(*types.Named); ok && i.Obj().Name() == "error" && i.Obj().Pkg() == nil {
			returnError = true
		} else {
			return fmt.Errorf("second return parameter must have type error but had: %s", result.At(1).Type())
		}
	}

	builderProperties := map[string]*builder.FieldMapping{}

	for name, property := range methodComments.Fields {
		builderProperties[name] = &builder.FieldMapping{
			Ignore: property.Ignore,
			Source: property.Source,
		}
		if property.Function != "" {
			def, err := g.parseMapExtend(converter, scope, property.Function)
			if err != nil {
				return err
			}
			builderProperties[name].Function = def
		}
	}

	m := &methodDefinition{
		Call:                   jen.Id(xtype.ThisVar).Dot(methodType.Name()),
		ID:                     methodType.String(),
		Explicit:               true,
		Name:                   methodType.Name(),
		Source:                 xtype.TypeOf(source),
		Target:                 xtype.TypeOf(target),
		Fields:                 builderProperties,
		MatchIgnoreCase:        methodComments.MatchIgnoreCase,
		IgnoreUnexportedFields: g.ignoreUnexportedFields,
		WrapErrors:             methodComments.WrapErrors,
		ReturnError:            returnError,
		ReturnTypeOrigin:       methodType.FullName(),
	}

	g.lookup[xtype.SignatureOf(m.Source, m.Target)] = m
	g.namer.Register(m.Name)
	return nil
}

func (g *generator) validateMethods() error {
	for _, method := range g.lookup {
		if method.Explicit && len(method.Fields) > 0 {
			isTargetStructPointer := method.Target.Pointer && method.Target.PointerInner.Struct
			if !method.Target.Struct && !isTargetStructPointer {
				return fmt.Errorf("Invalid struct field mapping on method:\n    %s\n\nField mappings like goverter:map or goverter:ignore may only be set on struct or struct pointers.\nSee https://goverter.jmattheis.de/#/conversion/configure-nested", method.ID)
			}
			if overlapping, ok := g.findOverlappingExplicitStructMethod(method); ok {
				return fmt.Errorf("Overlapping field mapping found.\n\nPlease move the field related goverter:* comments from this method:\n    %s\n\nto this method:\n    %s\n\nGoverter will use %s inside the implementation of %s, thus, field related settings would be ignored.", method.ID, overlapping.ID, overlapping.Name, method.Name)
			}
		}
	}
	return nil
}

func (g *generator) createMethods() error {
	methods := []*methodDefinition{}
	for _, method := range g.lookup {
		methods = append(methods, method)
	}
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name < methods[j].Name
	})
	for _, method := range methods {
		if method.Jen != nil && !method.Dirty {
			continue
		}
		method.Dirty = false
		err := g.buildMethod(method, builder.NoWrap)
		if err != nil {
			err = err.Lift(&builder.Path{
				SourceID:   "source",
				TargetID:   "target",
				SourceType: method.Source.T.String(),
				TargetType: method.Target.T.String(),
			})
			return fmt.Errorf("Error while creating converter method:\n    %s\n\n%s", method.ID, builder.ToString(err))
		}
	}
	for _, method := range g.lookup {
		if method.Dirty {
			return g.createMethods()
		}
	}
	g.appendToFile()
	return nil
}

func (g *generator) appendToFile() {
	methods := []*methodDefinition{}
	for _, method := range g.lookup {
		methods = append(methods, method)
	}
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name < methods[j].Name
	})
	for _, method := range methods {
		g.file.Add(method.Jen)
	}
}

func (g *generator) findOverlappingExplicitStructMethod(method *methodDefinition) (*methodDefinition, bool) {
	source := method.Source
	target := method.Target

	switch {
	case source.Struct && target.Pointer && target.PointerInner.Struct:
		method, ok := g.lookup[xtype.SignatureOf(source, target.PointerInner)]
		if ok && method.Explicit {
			return method, true
		}
	case source.Pointer && source.PointerInner.Struct && target.Pointer && target.PointerInner.Struct:
		method, ok := g.lookup[xtype.SignatureOf(source.PointerInner, target)]
		if ok && method.Explicit {
			return method, true
		}
		method, ok = g.lookup[xtype.SignatureOf(source.PointerInner, target.PointerInner)]
		if ok && method.Explicit {
			return method, true
		}
	}
	return nil, false
}

func (g *generator) buildMethod(method *methodDefinition, errWrapper builder.ErrorWrapper) *builder.Error {
	sourceID := jen.Id("source")
	source := method.Source
	target := method.Target

	returns := []jen.Code{target.TypeAsJen()}
	if method.ReturnError {
		returns = append(returns, jen.Id("error"))
	}

	fieldsTarget := method.Target.T.String()
	if method.Target.Pointer && method.Target.PointerInner.Struct {
		fieldsTarget = method.Target.PointerInner.T.String()
	}

	ctx := &builder.MethodContext{
		Namer:                  namer.New(),
		Fields:                 method.Fields,
		FieldsTarget:           fieldsTarget,
		IgnoreUnexportedFields: method.IgnoreUnexportedFields,
		MatchIgnoreCase:        method.MatchIgnoreCase,
		WrapErrors:             method.WrapErrors,
		TargetType:             method.Target,
		Signature:              xtype.SignatureOf(method.Source, method.Target),
	}

	var funcBlock []jen.Code
	if extendMethod, ok := g.extend[ctx.Signature]; ok {
		jenReturn, err := g.delegateMethod(
			ctx, extendMethod, xtype.VariableID(sourceID.Clone()), source, target, errWrapper)
		if err != nil {
			return err
		}
		funcBlock = []jen.Code{jenReturn}
	} else {
		stmt, newID, err := g.buildNoLookup(ctx, xtype.VariableID(sourceID.Clone()), source, target)
		if err != nil {
			return err
		}
		ret := []jen.Code{newID.Code}
		if method.ReturnError {
			ret = append(ret, jen.Nil())
		}

		funcBlock = append(stmt, jen.Return(ret...))
	}

	method.Jen = jen.Func().Params(jen.Id(xtype.ThisVar).Op("*").Id(g.name)).Id(method.Name).
		Params(jen.Id("source").Add(source.TypeAsJen())).Params(returns...).
		Block(funcBlock...)

	return nil
}

func (g *generator) buildNoLookup(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *builder.Error) {
	for _, rule := range BuildSteps {
		if rule.Matches(source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}
	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
}

type callableMethod struct {
	ID               string
	Call             *jen.Statement
	SelfAsFirstParam bool
	ReturnError      bool
	ReturnTypeOrigin string
}

func (g *generator) callByDefinition(
	ctx *builder.MethodContext,
	method *methodDefinition,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	callable := &callableMethod{
		ID:               method.ID,
		Call:             method.Call,
		ReturnError:      method.ReturnError,
		ReturnTypeOrigin: method.ReturnTypeOrigin,
		SelfAsFirstParam: method.SelfAsFirstParam,
	}
	return g.callMethod(ctx, callable, sourceID, source, target, errWrapper)
}

func (g *generator) CallExtendMethod(
	ctx *builder.MethodContext,
	method *builder.ExtendMethod,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	callable := &callableMethod{
		ID:               method.ID,
		Call:             method.Call,
		ReturnError:      method.ReturnError,
		ReturnTypeOrigin: method.ID,
		SelfAsFirstParam: method.SelfAsFirstParam,
	}
	return g.callMethod(ctx, callable, sourceID, source, target, errWrapper)
}

func (g *generator) callMethod(
	ctx *builder.MethodContext,
	method *callableMethod,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	params := []jen.Code{}
	if method.SelfAsFirstParam {
		params = append(params, jen.Id(xtype.ThisVar))
	}
	if sourceID != nil {
		params = append(params, sourceID.Code.Clone())
	}
	if method.ReturnError {
		current := g.lookup[ctx.Signature]
		if !current.ReturnError {
			if current.Explicit {
				return nil, nil, builder.NewError(fmt.Sprintf("ReturnTypeMismatch: Cannot use\n\n    %s\n\nin\n\n    %s\n\nbecause no error is returned as second return parameter", method.ReturnTypeOrigin, current.ID))
			}
			current.ReturnError = true
			current.ReturnTypeOrigin = method.ID
			current.Dirty = true
		}

		var errBlock []jen.Code
		if ctx.TargetVar == nil {
			innerName := ctx.Name("errValue")
			errBlock = []jen.Code{
				jen.Var().Id(innerName).Add(ctx.TargetType.TypeAsJen()),
				jen.Return(jen.Id(innerName), g.wrap(ctx, errWrapper, jen.Id("err"))),
			}
		} else {
			errBlock = []jen.Code{
				jen.Return(ctx.TargetVar, g.wrap(ctx, errWrapper, jen.Id("err"))),
			}
		}
		name := ctx.Name(target.ID())
		stmt := []jen.Code{
			jen.List(jen.Id(name), jen.Id("err")).Op(":=").Add(method.Call.Clone().Call(params...)),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(errBlock...),
		}
		return stmt, xtype.VariableID(jen.Id(name)), nil
	}
	id := xtype.OtherID(method.Call.Clone().Call(params...))
	return nil, id, nil
}

func (g *generator) delegateMethod(
	ctx *builder.MethodContext,
	delegateTo *methodDefinition,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) (*jen.Statement, *builder.Error) {
	params := []jen.Code{}
	if delegateTo.SelfAsFirstParam {
		params = append(params, jen.Id(xtype.ThisVar))
	}
	if sourceID != nil {
		params = append(params, sourceID.Code.Clone())
	}
	current := g.lookup[ctx.Signature]

	returns := []jen.Code{delegateTo.Call.Clone().Call(params...)}

	if delegateTo.ReturnError {
		if !current.ReturnError {
			return nil, builder.NewError(fmt.Sprintf("ReturnTypeMismatch: Cannot use\n\n    %s\n\nin\n\n    %s\n\nbecause no error is returned as second return parameter", delegateTo.ReturnTypeOrigin, current.ID))
		}
	} else {
		if current.ReturnError {
			returns = append(returns, jen.Nil())
		}
	}
	return jen.Return(returns...), nil
}

// wrap invokes the error wrapper if feature is enabled.
func (g *generator) wrap(ctx *builder.MethodContext, errWrapper builder.ErrorWrapper, errStmt *jen.Statement) *jen.Statement {
	if g.wrapErrors || ctx.WrapErrors {
		return errWrapper(errStmt)
	}
	return errStmt
}

// Build builds an implementation for the given source and target type, or uses an existing method for it.
func (g *generator) Build(
	ctx *builder.MethodContext,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	signature := xtype.SignatureOf(source, target)
	method, ok := g.extend[signature]
	if !ok {
		method, ok = g.lookup[signature]
	}

	if ok {
		return g.callByDefinition(ctx, method, sourceID, source, target, errWrapper)
	}

	hasPointerStructMethod := false
	if source.Struct && target.Struct {
		pointerSignature := xtype.SignatureOf(source.AsPointer(), target.AsPointer())

		_, hasPointerStructMethod = g.extend[pointerSignature]
		if !hasPointerStructMethod {
			_, hasPointerStructMethod = g.lookup[pointerSignature]
		}
	}

	if !hasPointerStructMethod && ((source.Named && !source.Basic) || (target.Named && !target.Basic) || (source.Pointer && source.PointerInner.Named && !source.PointerInner.Basic)) {
		name := g.namer.Name(source.UnescapedID() + "To" + strings.Title(target.UnescapedID()))

		method := &methodDefinition{
			ID:                     name,
			Name:                   name,
			Source:                 xtype.TypeOf(source.T),
			Target:                 xtype.TypeOf(target.T),
			Fields:                 map[string]*builder.FieldMapping{},
			IgnoreUnexportedFields: g.ignoreUnexportedFields,
			Call:                   jen.Id(xtype.ThisVar).Dot(name),
		}

		g.lookup[signature] = method
		g.namer.Register(method.Name)
		if err := g.buildMethod(method, errWrapper); err != nil {
			return nil, nil, err
		}
		// try again to trigger the found method thingy above
		return g.Build(ctx, sourceID, source, target, errWrapper)
	}

	for _, rule := range BuildSteps {
		if rule.Matches(source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}

	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
}

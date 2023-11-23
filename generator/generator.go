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

	Jen jen.Code
}

type generator struct {
	namer  *namer.Namer
	conf   *config.Converter
	lookup map[xtype.Signature]*generatedMethod
	extend map[xtype.Signature]*method.Definition
}

func (g *generator) buildMethods(f *jen.File) error {
	genMethods := []*generatedMethod{}
	for _, genMethod := range g.lookup {
		genMethods = append(genMethods, genMethod)
	}
	sort.Slice(genMethods, func(i, j int) bool {
		return genMethods[i].Name < genMethods[j].Name
	})
	for _, genMethod := range genMethods {
		if genMethod.Jen != nil && !genMethod.Dirty {
			continue
		}
		genMethod.Dirty = false
		err := g.buildMethod(genMethod, builder.NoWrap)
		if err != nil {
			err = err.Lift(&builder.Path{
				SourceID:   "source",
				TargetID:   "target",
				SourceType: genMethod.Source.String,
				TargetType: genMethod.Target.String,
			})
			return fmt.Errorf("Error while creating converter method:\n    %s\n\n%s", genMethod.ID, builder.ToString(err))
		}
	}
	for _, genMethod := range g.lookup {
		if genMethod.Dirty {
			return g.buildMethods(f)
		}
	}
	g.appendGenerated(f)
	return nil
}

func (g *generator) appendGenerated(f *jen.File) {
	genMethods := []*generatedMethod{}
	for _, genMethod := range g.lookup {
		genMethods = append(genMethods, genMethod)
	}
	sort.Slice(genMethods, func(i, j int) bool {
		return genMethods[i].Name < genMethods[j].Name
	})

	for _, def := range genMethods {
		f.Add(def.Jen)
	}
}

func (g *generator) buildMethod(genMethod *generatedMethod, errWrapper builder.ErrorWrapper) *builder.Error {
	sourceID := jen.Id("source")
	source := genMethod.Source
	target := genMethod.Target

	returns := []jen.Code{target.TypeAsJen()}
	if genMethod.ReturnError {
		returns = append(returns, jen.Id("error"))
	}

	fieldsTarget := genMethod.Target.String
	if genMethod.Target.Pointer && genMethod.Target.PointerInner.Struct {
		fieldsTarget = genMethod.Target.PointerInner.String
	}

	ctx := &builder.MethodContext{
		Namer:        namer.New(),
		Conf:         genMethod.Method,
		FieldsTarget: fieldsTarget,
		SeenNamed:    map[string]struct{}{},
		TargetType:   genMethod.Target,
		Signature:    genMethod.Signature(),
		HasMethod:    g.hasMethod,
	}

	var funcBlock []jen.Code
	if def, ok := g.extend[ctx.Signature]; ok {
		jenReturn, err := g.delegateMethod(
			ctx, def, xtype.VariableID(sourceID.Clone()), source, target, errWrapper)
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
		if genMethod.ReturnError {
			ret = append(ret, jen.Nil())
		}

		funcBlock = append(stmt, jen.Return(ret...))
	}

	genMethod.Jen = jen.Func().Params(jen.Id(xtype.ThisVar).Op("*").Id(g.conf.Name)).Id(genMethod.Name).
		Params(jen.Id("source").Add(source.TypeAsJen())).Params(returns...).
		Block(funcBlock...)

	return nil
}

func (g *generator) buildNoLookup(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *builder.Error) {
	for _, rule := range BuildSteps {
		if rule.Matches(ctx, source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}

	if source.Pointer && !target.Pointer {
		return nil, nil, builder.NewError(fmt.Sprintf(`TypeMismatch: Cannot convert %s to %s
It is unclear how nil should be handled in the pointer to non pointer conversion.

You can enable useZeroValueOnPointerInconsistency to instruct goverter to use the zero value if source is nil
https://goverter.jmattheis.de/#/config/useZeroValueOnPointerInconsistency

or you can define a custom conversion method with extend:
https://goverter.jmattheis.de/#/config/extend`, source.T, target.T))
	}

	return nil, nil, builder.NewError(fmt.Sprintf(`TypeMismatch: Cannot convert %s to %s

You can define a custom conversion method with extend:
https://goverter.jmattheis.de/#/config/extend`, source.T, target.T))
}

func (g *generator) CallMethod(
	ctx *builder.MethodContext,
	definition *method.Definition,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) ([]jen.Code, *xtype.JenID, *builder.Error) {
	params := []jen.Code{}
	if definition.SelfAsFirstParameter {
		params = append(params, jen.Id(xtype.ThisVar))
	}

	formatErr := func(s string) *builder.Error {
		return builder.NewError(fmt.Sprintf("Error using method:\n    %s\n\n%s", definition.ReturnTypeOriginID, s))
	}

	if definition.Source != nil {
		params = append(params, sourceID.Code.Clone())

		if definition.Source.String != source.String {
			cause := fmt.Sprintf("Method source type mismatches with conversion source: %s != %s", definition.Source.String, source.String)
			return nil, nil, formatErr(cause)
		}
	}

	if definition.Target.String != target.String {
		cause := fmt.Sprintf("Method return type mismatches with target: %s != %s", definition.Target.String, target.String)
		return nil, nil, formatErr(cause)
	}

	if definition.ReturnError {
		current := g.lookup[ctx.Signature]
		if !current.ReturnError {
			if current.Explicit {
				return nil, nil, formatErr("Used method returns error but conversion method does not")
			}
			current.ReturnError = true
			current.ReturnTypeOriginID = definition.ID
			current.Dirty = true
		}

		name := ctx.Name(target.ID())
		ctx.SetErrorTargetVar(jen.Id(name))

		stmt := []jen.Code{
			jen.List(jen.Id(name), jen.Id("err")).Op(":=").Add(definition.Call.Clone().Call(params...)),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(ctx.TargetVar, g.wrap(ctx, errWrapper, jen.Id("err")))),
		}
		return stmt, xtype.VariableID(jen.Id(name)), nil
	}
	id := xtype.OtherID(definition.Call.Clone().Call(params...))
	return nil, id, nil
}

func (g *generator) delegateMethod(
	ctx *builder.MethodContext,
	delegateTo *method.Definition,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	errWrapper builder.ErrorWrapper,
) (*jen.Statement, *builder.Error) {
	params := []jen.Code{}
	if delegateTo.SelfAsFirstParameter {
		params = append(params, jen.Id(xtype.ThisVar))
	}
	if sourceID != nil {
		params = append(params, sourceID.Code.Clone())
	}
	current := g.lookup[ctx.Signature]

	returns := []jen.Code{delegateTo.Call.Clone().Call(params...)}

	if delegateTo.ReturnError {
		if !current.ReturnError {
			return nil, builder.NewError(fmt.Sprintf("ReturnTypeMismatch: Cannot use\n\n    %s\n\nin\n\n    %s\n\nbecause no error is returned as second return parameter", delegateTo.ReturnTypeOriginID, current.ID))
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
	if ctx.Conf.WrapErrors {
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
	if def, ok := g.extend[signature]; ok {
		return g.CallMethod(ctx, def, sourceID, source, target, errWrapper)
	}
	if genMethod, ok := g.lookup[signature]; ok {
		return g.CallMethod(ctx, genMethod.Definition, sourceID, source, target, errWrapper)
	}

	isCurrentPointerStructMethod := false
	if source.Struct && target.Struct {
		// This checks if we are currently inside the generation of one of the following combinations.
		// *Source -> Target
		//  Source -> *Target
		// *Source -> *Target
		isCurrentPointerStructMethod = ctx.Signature.Source == source.AsPointerType().String() ||
			ctx.Signature.Target == target.AsPointerType().String()
	}

	createSubMethod := false

	if ctx.HasSeen(source) {
		g.lookup[ctx.Signature].Dirty = true
		createSubMethod = true
	} else if !isCurrentPointerStructMethod {
		switch {
		case source.Named && !source.Basic:
			createSubMethod = true
		case target.Named && !target.Basic:
			createSubMethod = true
		case source.Pointer && source.PointerInner.Named && !source.PointerInner.Basic:
			createSubMethod = true
		}
		if ctx.Conf.SkipCopySameType && source.String == target.String {
			createSubMethod = false
		}
	}
	ctx.MarkSeen(source)

	if createSubMethod {
		return g.createSubMethod(ctx, sourceID, source, target, errWrapper)
	}

	return g.buildNoLookup(ctx, sourceID, source, target)
}

func (g *generator) createSubMethod(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errWrapper builder.ErrorWrapper) ([]jen.Code, *xtype.JenID, *builder.Error) {
	if def, ok := g.getOverlappingStructDefinition(source, target); ok {
		return nil, nil, builder.NewError(fmt.Sprintf(`Overlapping struct settings found.

Move these field related settings:
    goverter:%s

from this method:
    %s

To a method you have to create with the following signature:
    func(%s) %s

Goverter will not use %s inside the current conversion method, thus, field related settings would be ignored.`, strings.Join(def.RawFieldSettings, "\n    goverter:"), def.ID, source.String, target.String, def.Name))
	}

	signature := xtype.SignatureOf(source, target)

	name := g.namer.Name(source.UnescapedID() + "To" + strings.Title(target.UnescapedID()))

	genMethod := &generatedMethod{
		Method: &config.Method{
			Common: g.conf.Common,
			Definition: &method.Definition{
				ID:   name,
				Name: name,
				Call: jen.Id(xtype.ThisVar).Dot(name),
				Parameters: method.Parameters{
					Source: xtype.TypeOf(source.T),
					Target: xtype.TypeOf(target.T),
				},
			},
		},
	}

	g.lookup[signature] = genMethod
	if err := g.buildMethod(genMethod, errWrapper); err != nil {
		return nil, nil, err
	}
	return g.CallMethod(ctx, genMethod.Definition, sourceID, source, target, errWrapper)
}

func (g *generator) hasMethod(source, target types.Type) bool {
	signature := xtype.Signature{Source: source.String(), Target: target.String()}
	_, ok := g.extend[signature]
	if !ok {
		_, ok = g.lookup[signature]
	}
	return ok
}

func (g *generator) getOverlappingStructDefinition(source, target *xtype.Type) (*generatedMethod, bool) {
	if !source.Struct || !target.Struct {
		return nil, false
	}

	overlapping := []xtype.Signature{
		{Source: source.AsPointerType().String(), Target: target.String},
		{Source: source.AsPointerType().String(), Target: target.AsPointerType().String()},
		{Source: source.String, Target: target.AsPointerType().String()},
	}

	for _, sig := range overlapping {
		if def, ok := g.lookup[sig]; ok && len(def.RawFieldSettings) > 0 {
			return def, true
		}
	}
	return nil, false
}

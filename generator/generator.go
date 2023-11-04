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
				SourceType: genMethod.Source.T.String(),
				TargetType: genMethod.Target.T.String(),
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

	fieldsTarget := genMethod.Target.T.String()
	if genMethod.Target.Pointer && genMethod.Target.PointerInner.Struct {
		fieldsTarget = genMethod.Target.PointerInner.T.String()
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
	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
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
	if definition.Source != nil {
		params = append(params, sourceID.Code.Clone())

		if definition.Source.T.String() != source.T.String() {
			cause := fmt.Sprintf("Method source type mismatches with conversion source: %s != %s", definition.Source.T.String(), source.T.String())
			return nil, nil, builder.NewError(cause).Lift(&builder.Path{
				Prefix:     "(",
				SourceID:   "source)",
				SourceType: definition.Source.T.String(),
			}).Lift(&builder.Path{
				Prefix:     ":",
				SourceID:   definition.Name,
				SourceType: definition.ID,
			})
		}
	}

	if definition.Target.T.String() != target.T.String() {
		cause := fmt.Sprintf("Method return type mismatches with target: %s != %s", definition.Target.T.String(), target.T.String())
		return nil, nil, builder.NewError(cause).Lift(&builder.Path{
			Prefix:     "(",
			SourceID:   ")",
			SourceType: definition.Parameters.Target.T.String(),
		}).Lift(&builder.Path{
			Prefix:     ":",
			SourceID:   definition.Name,
			SourceType: definition.ID,
		})
	}

	if definition.ReturnError {
		current := g.lookup[ctx.Signature]
		if !current.ReturnError {
			if current.Explicit {
				return nil, nil, builder.NewError(fmt.Sprintf("ReturnTypeMismatch: Cannot use\n\n    %s\n\nin\n\n    %s\n\nbecause no error is returned as second return parameter", definition.ReturnTypeOriginID, current.ID))
			}
			current.ReturnError = true
			current.ReturnTypeOriginID = definition.ID
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
			jen.List(jen.Id(name), jen.Id("err")).Op(":=").Add(definition.Call.Clone().Call(params...)),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(errBlock...),
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

	hasPointerStructMethod := false
	if source.Struct && target.Struct {
		pointerSignature := xtype.SignatureOf(source.AsPointer(), target.AsPointer())

		_, hasPointerStructMethod = g.extend[pointerSignature]
		if !hasPointerStructMethod {
			_, hasPointerStructMethod = g.lookup[pointerSignature]
		}
	}

	createSubMethod := false

	if ctx.HasSeen(source) {
		g.lookup[ctx.Signature].Dirty = true
		createSubMethod = true
	} else if !hasPointerStructMethod {
		switch {
		case source.Named && !source.Basic:
			createSubMethod = true
		case target.Named && !target.Basic:
			createSubMethod = true
		case source.Pointer && source.PointerInner.Named && !source.PointerInner.Basic:
			createSubMethod = true
		}
		if ctx.Conf.SkipCopySameType && source.T.String() == target.T.String() {
			createSubMethod = false
		}
	}
	ctx.MarkSeen(source)

	if createSubMethod {
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

	return g.buildNoLookup(ctx, sourceID, source, target)
}

func (g *generator) hasMethod(source, target types.Type) bool {
	signature := xtype.Signature{Source: source.String(), Target: target.String()}
	_, ok := g.extend[signature]
	if !ok {
		_, ok = g.lookup[signature]
	}
	return ok
}

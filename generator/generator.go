package generator

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/go-genconv/builder"
	"github.com/jmattheis/go-genconv/comments"
	"github.com/jmattheis/go-genconv/namer"
)

type Method struct {
	ID               string
	Immutable        bool
	Name             string
	Source           *builder.Type
	Target           *builder.Type
	AdditionalReturn []*builder.Type
	Mapping          map[string]string
	IgnoredFields    map[string]struct{}
	Delegate         *types.Func

	Jen jen.Code

	ReturnError bool
	Dirty       bool
}

type Generator struct {
	namer  *namer.Namer
	name   string
	file   *jen.File
	lookup map[builder.Signature]*Method
}

func (g *Generator) registerMethod(sources *types.Package, method *types.Func, methodComments comments.Method) error {
	signature, ok := method.Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("expected signature %#v", method.Type())
	}
	params := signature.Params()
	if params.Len() != 1 {
		return fmt.Errorf("expected signature to have only one parameter")
	}
	result := signature.Results()
	if result.Len() < 1 {
		return fmt.Errorf("return type must have at least one parameter")
	}
	if result.Len() > 2 {
		return fmt.Errorf("return type must have at max two parameter (target struct and error)")
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

	m := &Method{
		Immutable:     true,
		Name:          method.Name(),
		Source:        builder.TypeOf(source),
		Target:        builder.TypeOf(target),
		Mapping:       methodComments.NameMapping,
		IgnoredFields: methodComments.IgnoredFields,
		ReturnError:   returnError,
		ID:               method.String(),
	}

	if methodComments.Delegate != "" {
		delegate := sources.Scope().Lookup(methodComments.Delegate)
		if delegate == nil {
			return fmt.Errorf("delegate %s does not exist", methodComments.Delegate)
		}
		if f, ok := delegate.(*types.Func); ok {
			m.Delegate = f
		} else {
			return fmt.Errorf("delegate %s does is not a function", methodComments.Delegate)
		}
	}

	g.lookup[builder.Signature{
		Source: source.String(),
		Target: target.String(),
	}] = m
	g.namer.Register(m.Name)
	return nil
}

func (g *Generator) createMethods() error {
	methods := []*Method{}
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
		err := g.buildMethod(method)
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
	g.AppendToFile()
	return nil
}

func (g *Generator) AppendToFile() {
	methods := []*Method{}
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

func (g *Generator) buildMethod(method *Method) *builder.Error {
	sourceID := jen.Id("source")
	source := method.Source
	target := method.Target

	returns := []jen.Code{target.TypeAsJen()}
	if method.ReturnError {
		returns = append(returns, jen.Id("error"))
	}

	if method.Delegate != nil {
		method.Jen = jen.Func().Params(jen.Id("c").Op("*").Id(g.name)).Id(method.Name).
			Params(jen.Id("source").Add(source.TypeAsJen())).Params(returns...).
			Block(jen.Return(jen.Qual(method.Delegate.Pkg().Path(), method.Delegate.Name())).Call(jen.Id("c"), sourceID))
		return nil
	}

	ctx := &builder.MethodContext{
		Namer:         namer.New(),
		MappingBaseID: target.T.String(),
		Mapping:       method.Mapping,
		IgnoredFields: method.IgnoredFields,
		Signature:     builder.Signature{Source: method.Source.T.String(), Target: method.Target.T.String()},
	}
	stmt, newID, err := g.BuildNoLookup(ctx, builder.VariableID(sourceID.Clone()), source, target)
	if err != nil {
		return err
	}

	ret := []jen.Code{newID.Code}
	if method.ReturnError {
		ret = append(ret, jen.Nil())
	}

	stmt = append(stmt, jen.Return(ret...))

	method.Jen = jen.Func().Params(jen.Id("c").Op("*").Id(g.name)).Id(method.Name).
		Params(jen.Id("source").Add(source.TypeAsJen())).Params(returns...).
		Block(stmt...)

	return nil
}

func (g *Generator) BuildNoLookup(ctx *builder.MethodContext, sourceID *builder.JenID, source, target *builder.Type) ([]jen.Code, *builder.JenID, *builder.Error) {
	for _, rule := range BuildSteps {
		if rule.Matches(source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}
	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
}

func (g *Generator) Build(ctx *builder.MethodContext, sourceID *builder.JenID, source, target *builder.Type) ([]jen.Code, *builder.JenID, *builder.Error) {
	if method, ok := g.lookup[builder.Signature{Source: source.T.String(), Target: target.T.String()}]; ok {
		if method.ReturnError {
			current := g.lookup[ctx.Signature]
			if !current.ReturnError {
				current.ReturnError = true
				current.Dirty = true
			}

			name := ctx.Name(target.ID())
			stmt := []jen.Code{
				jen.List(jen.Id(name), jen.Id("err")).Op(":=").Id("c").Dot(method.Name).Call(sourceID.Code),
				jen.If(jen.Id("err").Op("!=").Nil()).Block(jen.Return(jen.Id(ctx.Namer.First), jen.Id("err"))),
			}
			return stmt, builder.VariableID(jen.Id(name)), nil
		}
		id := builder.OtherID(jen.Id("c").Dot(method.Name).Call(sourceID.Code))
		return nil, id, nil
	}

	if (source.Named && !source.Basic) || (target.Named && !target.Basic) {
		name := g.namer.Name(source.UnescapedID() + "To" + strings.Title(target.UnescapedID()))

		method := &Method{
			ID:            name,
			Name:          name,
			Source:        builder.TypeOf(source.T),
			Target:        builder.TypeOf(target.T),
			Mapping:       map[string]string{},
			IgnoredFields: map[string]struct{}{},
		}
		g.lookup[builder.Signature{Source: source.T.String(), Target: target.T.String()}] = method
		g.namer.Register(method.Name)
		if err := g.buildMethod(method); err != nil {
			return nil, nil, err
		}
		// try again to trigger the found method thingy above
		return g.Build(ctx, sourceID, source, target)
	}

	for _, rule := range BuildSteps {
		if rule.Matches(source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}

	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
}

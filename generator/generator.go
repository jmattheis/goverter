package generator

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/jmattheis/go-genconv/xtype"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/go-genconv/builder"
	"github.com/jmattheis/go-genconv/comments"
	"github.com/jmattheis/go-genconv/namer"
)

type Method struct {
	ID            string
	Explicit      bool
	Name          string
	Call          *jen.Statement
	Source        *xtype.Type
	Target        *xtype.Type
	Mapping       map[string]string
	IgnoredFields map[string]struct{}

	Jen jen.Code

	SelfAsFirstParam bool
	ReturnError      bool
	ReturnTypeOrigin string
	Dirty            bool
}

type Generator struct {
	namer  *namer.Namer
	name   string
	file   *jen.File
	lookup map[xtype.Signature]*Method
	extend map[xtype.Signature]*Method
}

func (g *Generator) registerMethod(method *types.Func, methodComments comments.Method) error {
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
		Call:             jen.Id(xtype.ThisVar).Dot(method.Name()),
		ID:               method.String(),
		Explicit:         true,
		Name:             method.Name(),
		Source:           xtype.TypeOf(source),
		Target:           xtype.TypeOf(target),
		Mapping:          methodComments.NameMapping,
		IgnoredFields:    methodComments.IgnoredFields,
		ReturnError:      returnError,
		ReturnTypeOrigin: method.FullName(),
	}

	g.lookup[xtype.Signature{
		Source: source.String(),
		Target: target.String(),
	}] = m
	g.namer.Register(m.Name)
	return nil
}

func (g *Generator) parseExtend(targetInterface types.Type, scope *types.Scope, methods []string) error {
	for _, method := range methods {
		obj := scope.Lookup(method)
		if obj == nil {
			return fmt.Errorf("%s does not exist in scope", method)
		}

		fn, ok := obj.(*types.Func)
		if !ok {
			return fmt.Errorf("%s is not a function", method)
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok {
			return fmt.Errorf("%s has no signature", method)
		}
		if sig.Params().Len() == 0 || sig.Results().Len() > 2 {
			return fmt.Errorf("%s has no or too many parameters", method)
		}
		if sig.Results().Len() == 0 || sig.Results().Len() > 2 {
			return fmt.Errorf("%s has no or too many returns", method)
		}

		source := sig.Params().At(0).Type()
		target := sig.Results().At(0).Type()
		returnError := false
		if sig.Results().Len() == 2 {
			if i, ok := sig.Results().At(1).Type().(*types.Named); ok && i.Obj().Name() == "error" && i.Obj().Pkg() == nil {
				returnError = true
			} else {
				return fmt.Errorf("second return parameter must have type error but had: %s", sig.Results().At(1).Type())
			}
		}

		selfAsFirstParameter := false
		if sig.Params().Len() == 2 {
			if source.String() == targetInterface.String() {
				selfAsFirstParameter = true
				source = sig.Params().At(1).Type()
			} else {
				return fmt.Errorf("the first parameter must be of type %s", targetInterface.String())
			}
		}

		g.extend[xtype.Signature{Source: source.String(), Target: target.String()}] = &Method{
			ID:               fn.String(),
			Explicit:         true,
			Call:             jen.Qual(fn.Pkg().Path(), fn.Name()),
			Name:             fn.Name(),
			Source:           xtype.TypeOf(source),
			Target:           xtype.TypeOf(target),
			SelfAsFirstParam: selfAsFirstParameter,
			ReturnError:      returnError,
			ReturnTypeOrigin: fn.String(),
		}

	}
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

	ctx := &builder.MethodContext{
		Namer:         namer.New(),
		Mapping:       method.Mapping,
		IgnoredFields: method.IgnoredFields,
		Signature:     xtype.Signature{Source: method.Source.T.String(), Target: method.Target.T.String()},
	}
	stmt, newID, err := g.BuildNoLookup(ctx, xtype.VariableID(sourceID.Clone()), source, target)
	if err != nil {
		return err
	}

	ret := []jen.Code{newID.Code}
	if method.ReturnError {
		ret = append(ret, jen.Nil())
	}

	stmt = append(stmt, jen.Return(ret...))

	method.Jen = jen.Func().Params(jen.Id(xtype.ThisVar).Op("*").Id(g.name)).Id(method.Name).
		Params(jen.Id("source").Add(source.TypeAsJen())).Params(returns...).
		Block(stmt...)

	return nil
}

func (g *Generator) BuildNoLookup(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *builder.Error) {
	for _, rule := range BuildSteps {
		if rule.Matches(source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}
	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
}

func (g *Generator) Build(ctx *builder.MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *builder.Error) {

	method, ok := g.extend[xtype.Signature{Source: source.T.String(), Target: target.T.String()}]
	if !ok {
		method, ok = g.lookup[xtype.Signature{Source: source.T.String(), Target: target.T.String()}]
	}

	if ok {
		params := []jen.Code{}
		if method.SelfAsFirstParam {
			params = append(params, jen.Id(xtype.ThisVar))
		}
		params = append(params, sourceID.Code.Clone())
		if method.ReturnError {
			current := g.lookup[ctx.Signature]
			if !current.ReturnError {
				if current.Explicit {
					return nil, nil, builder.NewError(fmt.Sprintf("ReturnTypeMismatch: Cannot use\n\n    %s\n\nin\n\n    %s\n\nbecause no error is returned as second parameter", method.ReturnTypeOrigin, current.ID))
				}
				current.ReturnError = true
				current.ReturnTypeOrigin = method.ID
				current.Dirty = true
			}

			name := ctx.Name(target.ID())
			stmt := []jen.Code{
				jen.List(jen.Id(name), jen.Id("err")).Op(":=").Add(method.Call.Clone().Call(params...)),
				jen.If(jen.Id("err").Op("!=").Nil()).Block(jen.Return(jen.Id(ctx.Namer.First), jen.Id("err"))),
			}
			return stmt, xtype.VariableID(jen.Id(name)), nil
		}
		id := xtype.OtherID(method.Call.Clone().Call(params...))
		return nil, id, nil
	}

	if (source.Named && !source.Basic) || (target.Named && !target.Basic) {
		name := g.namer.Name(source.UnescapedID() + "To" + strings.Title(target.UnescapedID()))

		method := &Method{
			ID:            name,
			Name:          name,
			Source:        xtype.TypeOf(source.T),
			Target:        xtype.TypeOf(target.T),
			Mapping:       map[string]string{},
			IgnoredFields: map[string]struct{}{},
			Call:          jen.Id(xtype.ThisVar).Dot(name),
		}
		g.lookup[xtype.Signature{Source: source.T.String(), Target: target.T.String()}] = method
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

package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// List handles array / slice types.
type List struct{}

// Matches returns true, if the builder can create handle the given types.
func (*List) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	if !source.List || !target.List {
		return false
	}

	if !target.ListFixed {
		return true
	}

	return source.ListFixed && target.ListFixed && source.ListLen == target.ListLen
}

// Build creates conversion source code for the given source and target type.
func (l *List) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	targetSlice := ctx.Name(target.ID())

	stmt, err := l.Assign(gen, ctx, AssignOf(jen.Id(targetSlice)), sourceID, source, target, path)
	if err != nil {
		return nil, nil, err
	}

	id := jen.Var().Add(jen.Id(targetSlice), target.TypeAsJen())
	stmt = append([]jen.Code{id}, stmt...)

	return stmt, xtype.VariableID(jen.Id(targetSlice)), nil
}

func (*List) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	index := ctx.Index()

	indexedSource := xtype.VariableID(sourceID.Code.Clone().Index(jen.Id(index)))

	forBlock, err := gen.Assign(ctx, assignTo.WithIndex(jen.Id(index)), indexedSource, source.ListInner, target.ListInner, path.Index(jen.Id(index)))
	if err != nil {
		return nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: source.ListInner.String,
			TargetID:   "[]",
			TargetType: target.ListInner.String,
		})
	}

	var result []jen.Code

	if !target.ListFixed {
		result = append(result, assignTo.Stmt.Clone().Op("=").Make(target.TypeAsJen(), jen.Len(sourceID.Code.Clone())))
	}

	forStmt := jen.For(jen.Id(index).Op(":=").Lit(0), jen.Id(index).Op("<").Len(sourceID.Code.Clone()), jen.Id(index).Op("++")).
		Block(forBlock...)
	result = append(result, forStmt)

	if source.ListFixed {
		return result, nil
	}

	return []jen.Code{jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(result...)}, nil
}

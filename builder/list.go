package builder

import "github.com/dave/jennifer/jen"

type List struct{}

func (*List) Matches(source, target *Type) bool {
	return source.List && target.List && !target.ListFixed
}

func (*List) Build(gen Generator, ctx *MethodContext, sourceID JenID, source, target *Type) ([]jen.Code, JenID, *Error) {
	targetSlice := ctx.Name("targetSlice")
	index := ctx.Index()

	indexedSource := sourceID.Clone().Index(jen.Id(index))

	newStmt, newId, err := gen.Build(ctx, indexedSource, source.ListInner, target.ListInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: source.ListInner.T.String(),
			TargetID:   "[]",
			TargetType: target.ListInner.T.String(),
		})
	}
	newStmt = append(newStmt, jen.Id(targetSlice).Index(jen.Id(index)).Op("=").Add(newId))

	stmt := []jen.Code{
		jen.Id(targetSlice).Op(":=").Make(target.TypeAsJen(), jen.Len(sourceID.Clone())),
		jen.For(jen.Id(index).Op(":=").Lit(0), jen.Id(index).Op("<").Len(sourceID.Clone()), jen.Id(index).Op("++")).
			Block(newStmt...),
	}

	return stmt, jen.Id(targetSlice), nil
}

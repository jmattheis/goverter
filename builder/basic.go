package builder

import "github.com/dave/jennifer/jen"

type Basic struct{}

func (*Basic) Matches(source, target *Type) bool {
	return source.Basic && target.Basic &&
		source.BasicType.Kind() == target.BasicType.Kind()
}

func (*Basic) Build(_ Generator, _ *MethodContext, sourceID JenID, source, target *Type) ([]jen.Code, JenID, *Error) {
	if target.Named || (!target.Named && source.Named) {
		return nil, target.TypeAsJen().Call(sourceID), nil
	}
	return nil, sourceID, nil
}

type BasicTargetPointerRule struct {
	basic Basic
}

func (*BasicTargetPointerRule) Matches(source, target *Type) bool {
	return source.Basic && target.Pointer && target.PointerInner.Basic
}

func (*BasicTargetPointerRule) Build(gen Generator, ctx *MethodContext, sourceID JenID, source, target *Type) ([]jen.Code, JenID, *Error) {
	name := ctx.Name("ref")

	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	stmt = append(stmt, jen.Id(name).Op(":=").Add(id))

	newId := jen.Op("&").Id(name)

	return stmt, newId, err
}

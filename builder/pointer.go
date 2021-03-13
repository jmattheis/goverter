package builder

import "github.com/dave/jennifer/jen"

type Pointer struct{}

func (*Pointer) Matches(source, target *Type) bool {
	return source.Pointer && target.Pointer
}

func (*Pointer) Build(gen Generator, ctx *MethodContext, sourceID JenID, source, target *Type) ([]jen.Code, JenID, *Error) {
	outerVar := ctx.Of(target, "outer")
	deref := ctx.Name("deref")

	nextBlock, id, err := gen.Build(ctx, jen.Id(deref), source.PointerInner, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}

	ifBlock := []jen.Code{
		jen.Id(deref).Op(":=").Op("*").Add(sourceID.Clone()),
	}
	ifBlock = append(ifBlock, nextBlock...)
	ifBlock = append(ifBlock, jen.Id(outerVar).Op("=").Op("&").Add(id))

	stmt := []jen.Code{
		jen.Var().Id(outerVar).Add(target.TypeAsJen()),
		jen.If(sourceID.Clone().Op("!=").Nil()).Block(ifBlock...),
	}
	return stmt, jen.Id(outerVar), err
}

type TargetPointer struct{}

func (*TargetPointer) Matches(source, target *Type) bool {
	return !source.Pointer && target.Pointer
}

func (*TargetPointer) Build(gen Generator, ctx *MethodContext, sourceID JenID, source, target *Type) ([]jen.Code, JenID, *Error) {
	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	newId := jen.Op("&").Add(id)
	return stmt, newId, nil
}

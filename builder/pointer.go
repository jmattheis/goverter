package builder

import "github.com/dave/jennifer/jen"

type Pointer struct{}

func (*Pointer) Matches(source, target *Type) bool {
	return source.Pointer && target.Pointer
}

func (*Pointer) Build(gen Generator, ctx *MethodContext, sourceID *JenID, source, target *Type) ([]jen.Code, *JenID, *Error) {
	outerVar := ctx.Name(target.ID())

	nextBlock, id, err := gen.Build(ctx, OtherID(jen.Op("*").Add(sourceID.Code.Clone())), source.PointerInner, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}

	ifBlock := nextBlock
	if id.Variable {
		ifBlock = append(ifBlock, jen.Id(outerVar).Op("=").Op("&").Add(id.Code.Clone()))
	} else {
		tempName := ctx.Name(target.PointerInner.ID())
		ifBlock = append(ifBlock, jen.Id(tempName).Op(":=").Add(id.Code.Clone()))
		ifBlock = append(ifBlock, jen.Id(outerVar).Op("=").Op("&").Id(tempName))
	}

	stmt := []jen.Code{
		jen.Var().Id(outerVar).Add(target.TypeAsJen()),
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(ifBlock...),
	}
	return stmt, VariableID(jen.Id(outerVar)), err
}

type TargetPointer struct{}

func (*TargetPointer) Matches(source, target *Type) bool {
	return !source.Pointer && target.Pointer
}

func (*TargetPointer) Build(gen Generator, ctx *MethodContext, sourceID *JenID, source, target *Type) ([]jen.Code, *JenID, *Error) {
	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	if id.Variable {
		return stmt, OtherID(jen.Op("&").Add(id.Code)), nil
	}
	tempName := ctx.Name(target.PointerInner.ID())
	stmt = append(stmt, jen.Id(tempName).Op(":=").Add(id.Code))
	return stmt, OtherID(jen.Op("&").Id(tempName)), nil
}

package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// Pointer handles pointer types.
type Pointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Pointer) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	return source.Pointer && target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*Pointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())

	stmt, outerVar, err := buildTargetVar(gen, ctx, sourceID, source, target)
	if err != nil {
		return nil, nil, err
	}

	ifBlock := []jen.Code{}

	valueSourceID := jen.Op("*").Add(sourceID.Code.Clone())
	if !source.PointerInner.Basic {
		valueSourceID = jen.Parens(valueSourceID)
	}

	innerID := xtype.OtherID(valueSourceID)
	innerID.ParentPointer = sourceID
	nextBlock, id, err := gen.Build(
		ctx, innerID, source.PointerInner, target.PointerInner, NoWrap)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.String,
			TargetID:   "*",
			TargetType: target.PointerInner.String,
		})
	}
	ifBlock = append(ifBlock, nextBlock...)
	if id.Variable {
		ifBlock = append(ifBlock, outerVar.Clone().Op("=").Op("&").Add(id.Code.Clone()))
	} else {
		tempName := ctx.Name(target.PointerInner.ID())
		ifBlock = append(ifBlock, jen.Id(tempName).Op(":=").Add(id.Code.Clone()))
		ifBlock = append(ifBlock, outerVar.Clone().Op("=").Op("&").Id(tempName))
	}

	stmt = append(stmt,
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(ifBlock...),
	)

	return stmt, xtype.VariableID(outerVar), err
}

// SourcePointer handles type were only the source is a pointer.
type SourcePointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*SourcePointer) Matches(ctx *MethodContext, source, target *xtype.Type) bool {
	return ctx.Conf.UseZeroValueOnPointerInconsistency && source.Pointer && !target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*SourcePointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	valueSourceID := jen.Op("*").Add(sourceID.Code.Clone())
	if !source.PointerInner.Basic {
		valueSourceID = jen.Parens(valueSourceID)
	}

	innerID := xtype.OtherID(valueSourceID)
	innerID.ParentPointer = sourceID

	stmt, valueVar, err := buildTargetVar(gen, ctx, sourceID, source, target)
	if err != nil {
		return nil, nil, err
	}

	nextInner, nextID, err := gen.Build(ctx, innerID, source.PointerInner, target, NoWrap)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.String,
		})
	}

	stmt = append(stmt,
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(
			append(nextInner, valueVar.Clone().Op("=").Add(nextID.Code))...,
		),
	)

	return stmt, xtype.VariableID(valueVar), nil
}

// TargetPointer handles type were only the target is a pointer.
type TargetPointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*TargetPointer) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	return !source.Pointer && target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*TargetPointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner, NoWrap)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.String,
			TargetID:   "*",
			TargetType: target.PointerInner.String,
		})
	}
	if id.Variable {
		return stmt, xtype.OtherID(jen.Op("&").Add(id.Code)), nil
	}
	tempName := ctx.Name(target.PointerInner.ID())
	stmt = append(stmt, jen.Id(tempName).Op(":=").Add(id.Code))
	return stmt, xtype.OtherID(jen.Op("&").Id(tempName)), nil
}

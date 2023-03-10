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
	outerVar := ctx.Name(target.ID())
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
			SourceType: source.PointerInner.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	ifBlock = append(ifBlock, nextBlock...)
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
	return stmt, xtype.VariableID(jen.Id(outerVar)), err
}

// SourcePointer handles type were only the source is a pointer.
type SourcePointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*SourcePointer) Matches(ctx *MethodContext, source, target *xtype.Type) bool {
	return ctx.Flags.Has(FlagZeroValueOnPtrInconsistency) && source.Pointer && !target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*SourcePointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Id(target.ID()))

	valueSourceID := jen.Op("*").Add(sourceID.Code.Clone())
	if !source.PointerInner.Basic {
		valueSourceID = jen.Parens(valueSourceID)
	}

	innerID := xtype.OtherID(valueSourceID)
	innerID.ParentPointer = sourceID

	valueVar := ctx.Name(target.ID())

	nextInner, nextID, err := gen.Build(ctx, innerID, source.PointerInner, target, NoWrap)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.T.String(),
		})
	}

	stmt := []jen.Code{
		jen.Var().Id(valueVar).Add(target.TypeAsJen()),
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(
			append(nextInner, jen.Id(valueVar).Op("=").Add(nextID.Code))...,
		),
	}

	return stmt, xtype.VariableID(jen.Id(valueVar)), nil
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
			SourceType: source.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	if id.Variable {
		return stmt, xtype.OtherID(jen.Op("&").Add(id.Code)), nil
	}
	tempName := ctx.Name(target.PointerInner.ID())
	stmt = append(stmt, jen.Id(tempName).Op(":=").Add(id.Code))
	return stmt, xtype.OtherID(jen.Op("&").Id(tempName)), nil
}

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
func (p *Pointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	return BuildByAssign(p, gen, ctx, sourceID, source, target, errPath)
}

func (*Pointer) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())

	valueSourceID := jen.Op("*").Add(sourceID.Code.Clone())
	if !source.PointerInner.Basic {
		valueSourceID = jen.Parens(valueSourceID)
	}

	innerID := xtype.OtherID(valueSourceID)
	innerID.ParentPointer = sourceID
	nextBlock, id, err := gen.Build(
		ctx, innerID, source.PointerInner, target.PointerInner, errPath)
	if err != nil {
		return nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.String,
			TargetID:   "*",
			TargetType: target.PointerInner.String,
		})
	}

	pstmt, tmpID := id.Pointer(target.PointerInner, ctx.Name)

	ifBlock := append(nextBlock, pstmt...)
	ifBlock = append(ifBlock, assignTo.Stmt.Clone().Op("=").Add(tmpID.Code))

	var elseCase []jen.Code
	if assignTo.Must {
		elseCase = append(elseCase, jen.Else().Block(assignTo.Stmt.Clone().Op("=").Nil()))
	}

	stmt := []jen.Code{
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(ifBlock...).Add(elseCase...),
	}

	return stmt, err
}

// SourcePointer handles type were only the source is a pointer.
type SourcePointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*SourcePointer) Matches(ctx *MethodContext, source, target *xtype.Type) bool {
	return ctx.Conf.UseZeroValueOnPointerInconsistency && source.Pointer && !target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (s *SourcePointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	return BuildByAssign(s, gen, ctx, sourceID, source, target, path)
}

func (*SourcePointer) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *Error) {
	valueSourceID := jen.Op("*").Add(sourceID.Code.Clone())
	if !source.PointerInner.Basic {
		valueSourceID = jen.Parens(valueSourceID)
	}

	innerID := xtype.OtherID(valueSourceID)
	innerID.ParentPointer = sourceID

	nextInner, nextID, err := gen.Build(ctx, innerID, source.PointerInner, target, path)
	if err != nil {
		return nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.String,
		})
	}

	stmt := []jen.Code{
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(
			append(nextInner, assignTo.Stmt.Clone().Op("=").Add(nextID.Code))...,
		),
	}

	return stmt, nil
}

// TargetPointer handles type were only the target is a pointer.
type TargetPointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*TargetPointer) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	return !source.Pointer && target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*TargetPointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner, path)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.String,
			TargetID:   "*",
			TargetType: target.PointerInner.String,
		})
	}

	pstmt, nextID := id.Pointer(target.PointerInner, ctx.Name)
	stmt = append(stmt, pstmt...)
	return stmt, nextID, nil
}

func (tp *TargetPointer) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *Error) {
	return AssignByBuild(tp, gen, ctx, assignTo, sourceID, source, target, path)
}

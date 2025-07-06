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
	if ctx.UseConstructor && ctx.Conf.DefaultUpdate {
		buildStmt, valueVar, err := buildTargetVar(gen, ctx, sourceID, source, target, errPath)
		if err != nil {
			return nil, nil, err
		}

		stmt, err := gen.Assign(ctx, AssignOf(jen.Parens(jen.Op("*").Add(valueVar))).IsUpdate(), sourceID.Deref(source), source.PointerInner, target.PointerInner, errPath)
		if err != nil {
			return nil, nil, err.Lift(&Path{
				SourceID:   "*",
				SourceType: source.PointerInner.String,
				TargetID:   "*",
				TargetType: target.PointerInner.String,
			})
		}

		buildStmt = append(buildStmt, jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(stmt...))

		return buildStmt, xtype.VariableID(valueVar), nil
	}

	return BuildByAssign(p, gen, ctx, sourceID, source, target, errPath)
}

func (*Pointer) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())

	nextBlock, id, err := gen.Build(ctx, sourceID.Deref(source), source.PointerInner, target.PointerInner, errPath)
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

	stmt := []jen.Code{
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(ifBlock...),
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
	if ctx.UseConstructor && ctx.Conf.DefaultUpdate {
		buildStmt, valueVar, err := buildTargetVar(gen, ctx, sourceID, source, target, path)
		if err != nil {
			return nil, nil, err
		}

		stmt, err := gen.Assign(ctx, AssignOf(valueVar).IsUpdate(), sourceID.Deref(source), source.PointerInner, target, path)
		if err != nil {
			return nil, nil, err.Lift(&Path{
				SourceID:   "*",
				SourceType: source.PointerInner.String,
			})
		}

		buildStmt = append(buildStmt, jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(stmt...))

		return buildStmt, xtype.VariableID(valueVar), nil
	}

	return BuildByAssign(s, gen, ctx, sourceID, source, target, path)
}

func (*SourcePointer) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, path ErrorPath) ([]jen.Code, *Error) {
	nextInner, err := gen.Assign(ctx, assignTo, sourceID.Deref(source), source.PointerInner, target, path)
	if err != nil {
		return nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.String,
		})
	}

	stmt := []jen.Code{
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).
			Block(nextInner...),
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

	if ctx.UseConstructor {
		buildStmt, valueVar, err := buildTargetVar(gen, ctx, sourceID, source, target, path)
		if err != nil {
			return nil, nil, err
		}

		stmt, err := gen.Assign(ctx, AssignOf(jen.Parens(jen.Op("*").Add(valueVar))).IsUpdate(), sourceID, source, target.PointerInner, path)
		if err != nil {
			return nil, nil, err.Lift(&Path{
				TargetID:   "*",
				TargetType: target.PointerInner.String,
			})
		}

		buildStmt = append(buildStmt, stmt...)

		return buildStmt, xtype.VariableID(valueVar), nil
	}

	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner, path)
	if err != nil {
		return nil, nil, err.Lift(&Path{
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

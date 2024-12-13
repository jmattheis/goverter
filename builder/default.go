package builder

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

func buildTargetVar(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *jen.Statement, *Error) {
	if !ctx.UseConstructor ||
		!types.Identical(ctx.Conf.Source.T, source.T) ||
		!types.Identical(ctx.Conf.Target.T, target.T) {
		name := ctx.Name(target.ID())
		variable := jen.Var().Id(name).Add(target.TypeAsJen())
		ctx.SetErrorTargetVar(jen.Id(name))
		return []jen.Code{variable}, jen.Id(name), nil
	}
	ctx.UseConstructor = false

	callTarget := target
	toPointer := target.Pointer && !ctx.Conf.Constructor.Target.Pointer
	if toPointer {
		callTarget = target.PointerInner
	}

	stmt, nextID, err := gen.CallMethod(ctx, ctx.Conf.Constructor, sourceID, source, callTarget, errPath)
	if err != nil {
		return nil, nil, err
	}

	if toPointer {
		pstmt, pointerID := nextID.Pointer(callTarget, ctx.Name)
		stmt = append(stmt, pstmt...)
		nextID = pointerID
	}

	if nextID.Variable {
		ctx.SetErrorTargetVar(nextID.Code.Clone())
		return stmt, nextID.Code, nil
	}
	name := ctx.Name(target.ID())
	stmt = append(stmt, jen.Id(name).Op(":=").Add(nextID.Code))
	ctx.SetErrorTargetVar(jen.Id(name))
	return stmt, jen.Id(name), nil
}

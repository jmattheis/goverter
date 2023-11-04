package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

func buildTargetVar(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *jen.Statement, *Error) {
	if ctx.Conf.Constructor == nil ||
		ctx.Conf.Source.T.String() != source.T.String() ||
		ctx.Conf.Target.T.String() != target.T.String() {
		name := ctx.Name(target.ID())
		variable := jen.Var().Id(name).Add(target.TypeAsJen())
		ctx.SetErrorTargetVar(jen.Id(name))
		return []jen.Code{variable}, jen.Id(name), nil
	}

	callTarget := target
	toPointer := target.Pointer && !ctx.Conf.Constructor.Target.Pointer
	if toPointer {
		callTarget = target.PointerInner
	}

	stmt, nextID, err := gen.CallMethod(ctx, ctx.Conf.Constructor, sourceID, source, callTarget, NoWrap)
	if err != nil {
		return nil, nil, err
	}

	if toPointer {
		variable := nextID.Code
		if !nextID.Variable {
			name := ctx.Name(target.ID() + "Val")
			stmt = append(stmt, jen.Id(name).Op(":=").Add(nextID.Code))
			variable = jen.Id(name)
		}
		nextID = xtype.OtherID(jen.Op("&").Add(variable))
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

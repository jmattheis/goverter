package xtype

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

func ZeroValue(t types.Type) *jen.Statement {
	switch cast := t.(type) {
	case *types.Basic:
		if cast.Info()&types.IsString != 0 {
			return jen.Lit("")
		} else if cast.Info()&types.IsNumeric != 0 {
			return jen.Lit(0)
		} else if cast.Info()&types.IsBoolean != 0 {
			return jen.Lit(false)
		}
		panic("unknown basic type" + cast.String())
	case *types.Named:
		switch under := cast.Underlying().(type) {
		case *types.Struct:
			return jen.Parens(toCode(t).Block())
		default:
			return ZeroValue(under)
		}
	case *types.Struct, *types.Array:
		return toCode(t).Block()
	case *types.Interface, *types.Signature, *types.Pointer, *types.Map, *types.Slice, *types.Chan:
		return jen.Nil()
	}
	panic("unsupported type " + t.String())
}

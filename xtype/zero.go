package xtype

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

// NotZeroValueCheck returns an expression that is true when sourceCode does not
// hold the zero value of type t.
//
// For most types this is `sourceCode != <zero>` (where <zero> is `nil` for
// pointers, slices, maps, channels, functions and interfaces). Some struct and
// array types are not comparable because they contain a slice, map or function,
// and Go rejects `!=` on them, producing uncompilable generated code (see
// #227). For those, fall back to `!reflect.ValueOf(sourceCode).IsZero()`, which
// works for any type.
func NotZeroValueCheck(sourceCode *jen.Statement, t types.Type) *jen.Statement {
	if zeroValueIsNil(t) || types.Comparable(t) {
		return sourceCode.Clone().Op("!=").Add(ZeroValue(t))
	}
	return jen.Op("!").Add(
		jen.Qual("reflect", "ValueOf").Call(sourceCode.Clone()).Dot("IsZero").Call(),
	)
}

// zeroValueIsNil reports whether the zero value of t is nil, i.e. t can be
// compared against nil regardless of whether it is otherwise comparable
// (slices and maps are not comparable but their zero value is nil).
func zeroValueIsNil(t types.Type) bool {
	switch t.Underlying().(type) {
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Signature, *types.Interface:
		return true
	}
	return false
}

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

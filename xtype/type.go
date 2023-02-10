package xtype

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
)

// ThisVar is used as name for the reference to the converter interface.
const ThisVar = "c"

// Signature represents a signature for conversion.
type Signature struct {
	Source string
	Target string
}

// Type is a helper wrapper for types.Type.
type Type struct {
	T             types.Type
	Interface     bool
	InterfaceType *types.Interface
	Struct        bool
	StructType    *types.Struct
	Named         bool
	NamedType     *types.Named
	Pointer       bool
	PointerType   *types.Pointer
	PointerInner  *Type
	List          bool
	ListFixed     bool
	ListInner     *Type
	Map           bool
	MapType       *types.Map
	MapKey        *Type
	MapValue      *Type
	Basic         bool
	BasicType     *types.Basic
}

// StructField holds the type of a struct field and its name.
type StructField struct {
	Name string
	Type *Type
}

// StructField returns the type of a struct field and its name upon successful match or
// an error if it is not found. This method will also return a detailed error if matchIgnoreCase
// is enabled and there are multiple non-exact matches.
func (t Type) StructField(name string, ignoreCase bool, ignored func(name string) bool) (*StructField, error) {
	if !t.Struct {
		panic("trying to get field of non struct")
	}

	var ambMatches []*StructField
	for y := 0; y < t.StructType.NumFields(); y++ {
		m := t.StructType.Field(y)
		if ignored(m.Name()) {
			continue
		}
		if m.Name() == name {
			// exact match takes precedence over case-insensitive match
			return &StructField{Name: m.Name(), Type: TypeOf(m.Type())}, nil
		}
		if ignoreCase && strings.EqualFold(m.Name(), name) {
			ambMatches = append(ambMatches, &StructField{Name: m.Name(), Type: TypeOf(m.Type())})
			// keep going to ensure struct does not have another case-insensitive match
		}
	}

	switch len(ambMatches) {
	case 0:
		return nil, fmt.Errorf("%q does not exist", name)
	case 1:
		return ambMatches[0], nil
	default:
		ambNames := make([]string, 0, len(ambMatches))
		for _, m := range ambMatches {
			ambNames = append(ambNames, m.Name)
		}
		return nil, ambiguousMatchError(name, ambNames)
	}
}

// JenID a jennifer code wrapper with extra infos.
type JenID struct {
	Code     *jen.Statement
	Variable bool
}

// VariableID is used, when the ID can be referenced. F.ex it is not a function call.
func VariableID(code *jen.Statement) *JenID {
	return &JenID{Code: code, Variable: true}
}

// OtherID is used, when the ID isn't a variable id.
func OtherID(code *jen.Statement) *JenID {
	return &JenID{Code: code, Variable: false}
}

// TypeOf creates a Type.
func TypeOf(t types.Type) *Type {
	rt := &Type{}
	rt.T = t
	switch value := t.(type) {
	case *types.Pointer:
		rt.Pointer = true
		rt.PointerType = value
		rt.PointerInner = TypeOf(value.Elem())
	case *types.Basic:
		rt.Basic = true
		rt.BasicType = value
	case *types.Map:
		rt.Map = true
		rt.MapType = value
		rt.MapKey = TypeOf(value.Key())
		rt.MapValue = TypeOf(value.Elem())
	case *types.Slice:
		rt.List = true
		rt.ListInner = TypeOf(value.Elem())
	case *types.Array:
		rt.List = true
		rt.ListFixed = true
		rt.ListInner = TypeOf(value.Elem())
	case *types.Named:
		underlying := TypeOf(value.Underlying())
		underlying.T = value
		underlying.Named = true
		underlying.NamedType = value
		return underlying
	case *types.Struct:
		rt.Struct = true
		rt.StructType = value
	case *types.Interface:
		rt.Interface = true
		rt.InterfaceType = value
	default:
		panic("unknown types.Type " + t.String())
	}
	return rt
}

// ID returns a deteministically generated id that may be used as variable.
func (t *Type) ID() string {
	return t.asID(true, true)
}

// UnescapedID returns a deteministically generated id that may be used as variable
// reserved keywords aren't escaped.
func (t *Type) UnescapedID() string {
	return t.asID(true, false)
}

func (t *Type) asID(seeNamed, escapeReserved bool) string {
	if seeNamed && t.Named {
		pkgName := t.NamedType.Obj().Pkg().Name()
		name := pkgName + t.NamedType.Obj().Name()
		return name
	}
	if t.List {
		return t.ListInner.asID(true, false) + "List"
	}
	if t.Basic {
		if escapeReserved {
			return "x" + t.BasicType.String()
		}
		return t.BasicType.String()
	}
	if t.Pointer {
		return "p" + strings.Title(t.PointerInner.asID(true, false))
	}
	if t.Map {
		return "map" + strings.Title(t.MapKey.asID(true, false)+strings.Title(t.MapValue.asID(true, false)))
	}
	if t.Struct {
		return "unnamed"
	}
	return "unknown"
}

// TypeAsJen returns a jen representation of the type.
func (t Type) TypeAsJen() *jen.Statement {
	if t.Named {
		return toCode(t.NamedType, &jen.Statement{})
	}
	return toCode(t.T, &jen.Statement{})
}

func toCode(t types.Type, st *jen.Statement) *jen.Statement {
	switch cast := t.(type) {
	case *types.Named:
		if cast.Obj().Pkg() == nil {
			return st.Id(cast.Obj().Name())
		}
		return st.Qual(cast.Obj().Pkg().Path(), cast.Obj().Name())
	case *types.Map:
		key := toCode(cast.Key(), &jen.Statement{})
		return toCode(cast.Elem(), st.Map(key))
	case *types.Slice:
		return toCode(cast.Elem(), st.Index())
	case *types.Array:
		return toCode(cast.Elem(), st.Index(jen.Lit(int(cast.Len()))))
	case *types.Pointer:
		return toCode(cast.Elem(), st.Op("*"))
	case *types.Basic:
		return toCodeBasic(cast.Kind(), st)
	case *types.Struct:
		return toCodeStruct(cast, st)
	}
	panic("unsupported type " + t.String())
}

func toCodeStruct(t *types.Struct, st *jen.Statement) *jen.Statement {
	fields := []jen.Code{}
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		tag := t.Tag(i)

		fieldType := toCode(f.Type(), &jen.Statement{})
		if tag != "" {
			fieldType = fieldType.Add(jen.Id("`" + tag + "`"))
		}

		if !f.Embedded() {
			fieldType = jen.Id(f.Name()).Add(fieldType)
		}

		fields = append(fields, fieldType)
	}

	return st.Struct(fields...)
}

func toCodeBasic(t types.BasicKind, st *jen.Statement) *jen.Statement {
	switch t {
	case types.String:
		return st.String()
	case types.Int:
		return st.Int()
	case types.Int8:
		return st.Int8()
	case types.Int16:
		return st.Int16()
	case types.Int32:
		return st.Int32()
	case types.Int64:
		return st.Int64()
	case types.Uint:
		return st.Uint()
	case types.Uint8:
		return st.Uint8()
	case types.Uint16:
		return st.Uint16()
	case types.Uint32:
		return st.Uint32()
	case types.Uint64:
		return st.Uint64()
	case types.Bool:
		return st.Bool()
	case types.Complex128:
		return st.Complex128()
	case types.Complex64:
		return st.Complex64()
	case types.Float32:
		return st.Float32()
	case types.Float64:
		return st.Float64()
	default:
		panic(fmt.Sprintf("unsupported type %d", t))
	}
}

func ambiguousMatchError(name string, ambNames []string) error {
	return fmt.Errorf(`multiple matches found for %q. Possible matches: %s.

Explicitly define the mapping via goverter:map. Example:

    goverter:map %s %s

See https://goverter.jmattheis.de/#/conversion/mapping`, name, strings.Join(ambNames, ", "), ambNames[0], name)
}

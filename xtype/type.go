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

func SignatureOf(source, target *Type) Signature {
	return Signature{Source: source.String, Target: target.String}
}

// Type is a helper wrapper for types.Type.
type Type struct {
	String        string
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
	Signature     bool
	SignatureType *types.Signature
}

func (t *Type) AsPointer() *Type {
	return TypeOf(types.NewPointer(t.T))
}

// StructField holds the type of a struct field and its name.
type StructField struct {
	Path []string
	Type *Type
}

type SimpleStructField struct {
	Name string
	Type *Type
}

// StructField returns the type of a struct field and its name upon successful match or
// an error if it is not found. This method will also return a detailed error if matchIgnoreCase
// is enabled and there are multiple non-exact matches.
func (t Type) findAllFields(path []string, name string, ignoreCase bool, ignored func(name string) bool) (*StructField, []*StructField) {
	if !t.Struct {
		panic("trying to get field of non struct")
	}

	var matches []*StructField
	for y := 0; y < t.StructType.NumFields(); y++ {
		m := t.StructType.Field(y)
		if ignored(m.Name()) {
			continue
		}
		if m.Name() == name {
			// exact match takes precedence over case-insensitive match
			newPath := append([]string{}, path...)
			newPath = append(newPath, m.Name())
			return &StructField{Path: newPath, Type: TypeOf(m.Type())}, nil
		}
		if ignoreCase && strings.EqualFold(m.Name(), name) {
			newPath := append([]string{}, path...)
			newPath = append(newPath, m.Name())
			matches = append(matches, &StructField{Path: newPath, Type: TypeOf(m.Type())})
			// keep going to ensure struct does not have another case-insensitive match
		}
	}

	return nil, matches
}

type FieldSources struct {
	Path []string
	Type *Type
}

func FindExactField(source *Type, name string) (*SimpleStructField, error) {
	exactMatch, _ := source.findAllFields(nil, name, false, func(name string) bool { return false })
	if exactMatch == nil {
		return nil, fmt.Errorf("%q does not exist", name)
	}
	return &SimpleStructField{Name: exactMatch.Path[0], Type: exactMatch.Type}, nil
}

type NoMatchError struct{ Field string }

func (err *NoMatchError) Error() string {
	return fmt.Sprintf("\"%s\" does not exist", err.Field)
}

func FindField(name string, ignoreCase bool, ignored func(name string) bool, source *Type, additionalFieldSources []FieldSources) (*StructField, error) {
	exactMatch, ignoreCaseMatches := source.findAllFields(nil, name, ignoreCase, ignored)
	var exactMatches []*StructField
	if exactMatch != nil {
		exactMatches = append(exactMatches, exactMatch)
	}

	for _, source := range additionalFieldSources {
		sourceExactMatch, sourceIgnoreCaseMatches := source.Type.findAllFields(source.Path, name, ignoreCase, ignored)
		if sourceExactMatch != nil {
			exactMatches = append(exactMatches, sourceExactMatch)
		}
		ignoreCaseMatches = append(ignoreCaseMatches, sourceIgnoreCaseMatches...)
	}

	matches := exactMatches
	if len(matches) == 0 {
		matches = ignoreCaseMatches
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return nil, &NoMatchError{Field: name}
	default:
		names := make([]string, 0, len(matches))
		for _, m := range matches {
			names = append(names, strings.Join(m.Path, "."))
		}
		return nil, ambiguousMatchError(name, names)
	}
}

// JenID a jennifer code wrapper with extra infos.
type JenID struct {
	ParentPointer *JenID
	Code          *jen.Statement
	Variable      bool
}

func (j *JenID) Pointer(t *Type, namer func(string) string) ([]jen.Code, *JenID) {
	if j.Variable {
		return nil, OtherID(jen.Op("&").Add(j.Code.Clone()))
	}

	name := namer(t.ID())
	stmt := []jen.Code{jen.Id(name).Op(":=").Add(j.Code.Clone())}
	return stmt, OtherID(jen.Op("&").Id(name))
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
	rt.String = t.String()
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
		underlying.T = rt.T
		underlying.String = rt.String
		underlying.Named = true
		underlying.NamedType = value
		return underlying
	case *types.Struct:
		rt.Struct = true
		rt.StructType = value
	case *types.Interface:
		rt.Interface = true
		rt.InterfaceType = value
	case *types.Signature:
		rt.Signature = true
		rt.SignatureType = value
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
		return toCode(t.NamedType)
	}
	return toCode(t.T)
}

func ambiguousMatchError(name string, ambNames []string) error {
	return fmt.Errorf(`multiple matches found for %q. Possible matches: %s.

Explicitly define the mapping via goverter:map. Example:

    goverter:map %s %s

See https://goverter.jmattheis.de/#/config/map`, name, strings.Join(ambNames, ", "), ambNames[0], name)
}

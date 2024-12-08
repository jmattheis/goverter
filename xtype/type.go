package xtype

import (
	"fmt"
	"go/types"
	"reflect"
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
	Func          bool
	FuncType      *types.Func
	Chan          bool
	ChanType      *types.Chan

	enum *Enum
}

func (t *Type) AssignableTo(other *Type) bool {
	return types.AssignableTo(t.T, other.T)
}

func (t *Type) AsPointer() *Type {
	return TypeOf(t.AsPointerType())
}

func (t *Type) AsPointerType() *types.Pointer {
	return types.NewPointer(t.T)
}

func (t *Type) inStruct(source *Type, field string) *Type {
	if t.Signature && source.Named {
		t.FuncType = types.NewFunc(-1, source.NamedType.Obj().Pkg(), field, t.SignatureType)
		t.Func = true
	}

	return t
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
func (t Type) findAllFields(path []string, name, targetTag, matchTag string, ignoreCase bool) (*StructField, []*StructField) {
	if !t.Struct {
		panic("trying to get field of non struct")
	}

	var matches []*StructField
	tagName := parseTagName(matchTag, targetTag)
	handle := func(inType types.Type, fieldNum int) *StructField {
		var obj types.Object
		switch inT := inType.(type) {
		case *types.Struct:
			obj = inT.Field(fieldNum)
			if len(matchTag) == 0 || len(tagName) == 0 {
				break
			}
			if parseTagName(matchTag, inT.Tag(fieldNum)) == tagName {
				return &StructField{Path: append(path, obj.Name()), Type: TypeOf(obj.Type()).inStruct(&t, obj.Name())}
			}
		case *types.Named:
			obj = inT.Method(fieldNum)
		default:
			return nil
		}

		exact := obj.Name() == name
		if exact || (ignoreCase && strings.EqualFold(obj.Name(), name)) {
			// exact match takes precedence over case-insensitive match
			f := &StructField{Path: append(path, obj.Name()), Type: TypeOf(obj.Type()).inStruct(&t, obj.Name())}
			if exact {
				return f
			}
			matches = append(matches, f)
		}
		return nil
	}

	for y := 0; y < t.StructType.NumFields(); y++ {
		if exact := handle(t.StructType, y); exact != nil {
			return exact, matches
		}
	}

	if t.Named {
		for y := 0; y < t.NamedType.NumMethods(); y++ {
			if exact := handle(t.NamedType, y); exact != nil {
				return exact, matches
			}
		}
	}

	return nil, matches
}

type FieldSources struct {
	Path []string
	Type *Type
}

func FindExactField(source *Type, name, targetTag, matchTag string) (*SimpleStructField, error) {
	exactMatch, _ := source.findAllFields(nil, name, targetTag, matchTag, false)
	if exactMatch == nil {
		return nil, fmt.Errorf("%q does not exist", name)
	}
	return &SimpleStructField{Name: exactMatch.Path[0], Type: exactMatch.Type}, nil
}

type NoMatchError struct{ Field string }

func (err *NoMatchError) Error() string {
	return fmt.Sprintf("\"%s\" does not exist", err.Field)
}

func FindField(name, targetTag, matchTag string, ignoreCase bool, source *Type, additionalFieldSources []FieldSources) (*StructField, error) {
	exactMatch, ignoreCaseMatches := source.findAllFields(nil, name, targetTag, matchTag, ignoreCase)
	var exactMatches []*StructField
	if exactMatch != nil {
		exactMatches = append(exactMatches, exactMatch)
	}

	for _, source := range additionalFieldSources {
		sourceExactMatch, sourceIgnoreCaseMatches := source.Type.findAllFields(source.Path, name, targetTag, matchTag, ignoreCase)
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

// protobufSubtagOrder specifies which protobuf tag sub-tags should be considered and in which order.
var protobufSubtagOrder = [2]string{"json=", "name="}

// parseTagName parses out the field name from a struct tag.
func parseTagName(matchTag, tag string) string {
	entry := reflect.StructTag(tag).Get(matchTag)

	// Protobuf tags are exciting, in that they /optionally/ have a json= section which takes
	// precedence over the name= section.
	if matchTag == "protobuf" {
		// Load the subtags we care about into a map (avoids O(2N)).
		entryMap := make(map[string]string, 2)
		for _, section := range strings.Split(entry, ",") {
			for _, prefix := range protobufSubtagOrder {
				if strings.HasPrefix(section, prefix) {
					entryMap[prefix] = strings.Split(section, "=")[1]
				}
			}
		}
		// Now that we have the full list of subtags, choose one in order of preference.
		for _, prefix := range protobufSubtagOrder {
			if x, ok := entryMap[prefix]; ok {
				return x
			}
		}
		// Fall back to json.
		entry = reflect.StructTag(tag).Get("json")
	}

	// Standard tags.
	return strings.Split(entry, ",")[0]
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
	t = Unalias(t)
	rt := &Type{}
	rt.T = t
	rt.String = t.String()
	applyTo(rt, t)
	return rt
}

func applyTo(rt *Type, t types.Type) {
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
		rt.Named = true
		rt.NamedType = value
		applyTo(rt, value.Underlying())
	case *types.Struct:
		rt.Struct = true
		rt.StructType = value
	case *types.Interface:
		rt.Interface = true
		rt.InterfaceType = value
	case *types.Signature:
		rt.Signature = true
		rt.SignatureType = value
	case *types.Chan:
		rt.Chan = true
		rt.ChanType = value
	case *types.TypeParam:
		// ignore
	default:
		panic("unknown types.Type " + t.String())
	}
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
		pkg := t.NamedType.Obj().Pkg()
		name := t.NamedType.Obj().Name()
		switch {
		case pkg != nil:
			name = pkg.Name() + name
		case escapeReserved:
			name = "x" + name
		}
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
	if t.Chan {
		return "chan"
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

See https://goverter.jmattheis.de/reference/map`, name, strings.Join(ambNames, ", "), ambNames[0], name)
}

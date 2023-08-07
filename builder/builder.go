package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/namer"
	"github.com/jmattheis/goverter/xtype"
)

// Builder builds converter implementations, and can decide if it can handle the given type.
type Builder interface {
	// Matches returns true, if the builder can create handle the given types.
	Matches(ctx *MethodContext, source, target *xtype.Type) bool
	// Build creates conversion source code for the given source and target type.
	Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error)
}

// Generator checks all existing builders if they can create a conversion implementations for the given source and target type
// If no one Builder#Matches then, an error is returned.
type Generator interface {
	Build(
		ctx *MethodContext,
		sourceID *xtype.JenID,
		source, target *xtype.Type,
		errWrapper ErrorWrapper) ([]jen.Code, *xtype.JenID, *Error)

	CallExtendMethod(
		ctx *MethodContext,
		method *ExtendMethod,
		sourceID *xtype.JenID,
		source, target *xtype.Type,
		errWrapper ErrorWrapper) ([]jen.Code, *xtype.JenID, *Error)
}

// MethodContext exposes information for the current method.
type MethodContext struct {
	*namer.Namer
	Fields       map[string]*FieldMapping
	SubFields    map[string]*FieldMapping
	FieldsTarget string
	Signature    xtype.Signature
	TargetType   *xtype.Type
	AutoMap      []string
	Flags        ConversionFlags
	SeenNamed    map[string]struct{}

	TargetVar *jen.Statement
}

func (ctx *MethodContext) HasSeen(source *xtype.Type) bool {
	if !source.Named {
		return false
	}
	typeString := source.NamedType.String()
	_, ok := ctx.SeenNamed[typeString]
	return ok
}

func (ctx *MethodContext) MarkSeen(source *xtype.Type) {
	if !source.Named {
		return
	}
	typeString := source.NamedType.String()
	ctx.SeenNamed[typeString] = struct{}{}
}

func (ctx *MethodContext) SetErrorTargetVar(m *jen.Statement) {
	if ctx.TargetVar == nil {
		ctx.TargetVar = m
	}
}

func (ctx *MethodContext) Field(target *xtype.Type, name string) *FieldMapping {
	if ctx.FieldsTarget != target.T.String() {
		return emptyMapping
	}

	prop, ok := ctx.Fields[name]
	if !ok {
		return emptyMapping
	}
	return prop
}

var emptyMapping *FieldMapping = &FieldMapping{}

type FieldMapping struct {
	Source   string
	Ignore   bool
	Function *ExtendMethod
}
type ExtendMethod struct {
	ID               string
	Name             string
	SelfAsFirstParam bool
	ReturnError      bool
	Call             *jen.Statement
	// optional source
	Source *xtype.Type
	Target *xtype.Type
}

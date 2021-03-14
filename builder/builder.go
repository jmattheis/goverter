package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/go-genconv/namer"
	"github.com/jmattheis/go-genconv/xtype"
)

type Builder interface {
	Matches(source, target *xtype.Type) bool
	Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error)
}

type Generator interface {
	Build(ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error)
}

type MethodContext struct {
	*namer.Namer
	MappingBaseID string
	Mapping       map[string]string
	IgnoredFields map[string]struct{}
	Signature     xtype.Signature
}

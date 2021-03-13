package builder

import (
	"github.com/dave/jennifer/jen"
)

type Builder interface {
	Matches(source, target *Type) bool
	Build(gen Generator, ctx *MethodContext, sourceID *JenID, source, target *Type) ([]jen.Code, *JenID, *Error)
}

type Generator interface {
	Build(ctx *MethodContext, sourceID *JenID, source, target *Type) ([]jen.Code, *JenID, *Error)
}

type MethodContext struct {
	*Namer
	MappingBaseID string
	Mapping       map[string]string
	IgnoredFields map[string]struct{}
}

package method

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

type Definition struct {
	Parameters
	OriginID string
	Call     *jen.Statement
	ID       string
	Package  string
	Name     string

	Generated  bool
	CustomCall *jen.Statement
}

type Parameters struct {
	TypeParams bool

	Source       *xtype.Type
	MultiSources []*xtype.Type
	Target       *xtype.Type
	Context      map[string]*xtype.Type

	Signature xtype.Signature

	RawArgs []Arg

	ReturnError  bool
	UpdateTarget bool
}

type Arg struct {
	Name string
	Use  ArgUse
	Type *xtype.Type
}

type ArgUse string

const (
	ArgUseSource      ArgUse = "source"
	ArgUseMultiSource ArgUse = "additional-source"
	ArgUseInterface   ArgUse = "interface"
	ArgUseContext     ArgUse = "context"
	ArgUseTarget      ArgUse = "target"
)

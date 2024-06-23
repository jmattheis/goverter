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

func (m *Definition) Signature() xtype.Signature {
	return xtype.SignatureOf(m.Parameters.Source, m.Parameters.Target)
}

type Parameters struct {
	ReturnError          bool
	SelfAsFirstParameter bool
	TypeParams           bool
	Source               *xtype.Type
	Target               *xtype.Type
}

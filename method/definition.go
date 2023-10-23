package method

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

type Definition struct {
	Parameters
	ID                 string
	Call               *jen.Statement
	Name               string
	ReturnTypeOriginID string
}

func (m *Definition) Signature() xtype.Signature {
	return xtype.SignatureOf(m.Parameters.Source, m.Parameters.Target)
}

type Parameters struct {
	ReturnError          bool
	SelfAsFirstParameter bool
	Source               *xtype.Type
	Target               *xtype.Type
}

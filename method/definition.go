package method

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

type Definition struct {
	Parameters
	ID                 string
	Package            string
	Name               string
	ReturnTypeOriginID string

	Generated  bool
	CustomCall *jen.Statement
}

func (m *Definition) Qual(isStruct bool) *jen.Statement {
	switch {
	case m.CustomCall != nil:
		return m.CustomCall
	case isStruct && m.Generated:
		return jen.Id(xtype.ThisVar).Dot(m.Name)
	default:
		return jen.Qual(m.Package, m.Name)
	}
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

package identity

import (
	"github.com/jmattheis/goverter/xtype"
)

type Definition struct {
	OriginID string
	ID       string
	Package  string
	Name     string

	Type *xtype.Type
}

package identity

import (
	"go/types"
)

type Index struct {
	List []types.Type
}

func NewIndex() *Index {
	return &Index{}
}

func (l *Index) Has(t types.Type) bool {
	for _, tt := range l.List {
		if types.Identical(t, tt) {
			return true
		}
	}
	return false
}

func (l *Index) RegisterDefinition(def *Definition) {
	l.List = append(l.List, def.Type.T)
}

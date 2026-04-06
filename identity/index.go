package identity

import (
	"go/types"
)

type Index struct {
	List map[types.Type]struct{}
}

func NewIndex() *Index {
	return &Index{List: map[types.Type]struct{}{}}
}

func (l *Index) Has(t types.Type) bool {
	_, ok := l.List[t]
	return ok
}

func (l *Index) RegisterDefinition(def *Definition) {
	l.List[def.Type.T] = struct{}{}
}

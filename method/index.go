package method

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jmattheis/goverter/xtype"
)

type IndexEntry[T any] struct {
	Def  *Definition
	Item *T
}

type IndexID struct {
	sig xtype.Signature
	idx int
}

func NewIndex[T any]() *Index[T] {
	return &Index[T]{
		Exact: map[xtype.Signature][]IndexEntry[T]{},
	}
}

type Index[T any] struct {
	Exact map[xtype.Signature][]IndexEntry[T]
}

func (l *Index[T]) RegisterOverrideOverlapping(t *T, def *Definition) {
	newEntry := IndexEntry[T]{Def: def, Item: t}
	for i, entry := range l.Exact[def.Signature] {
		if satisfiesContext(entry.Def.Context, def.Context) || satisfiesContext(def.Context, entry.Def.Context) {
			l.Exact[def.Signature][i] = newEntry
			return
		}
	}

	l.Exact[def.Signature] = append(l.Exact[def.Signature], newEntry)
}

func (l *Index[T]) Register(t *T, def *Definition) (IndexID, error) {
	for _, entry := range l.Exact[def.Signature] {
		if err := checkOverlap(entry.Def, def); err != nil {
			return IndexID{}, err
		}
		if err := checkOverlap(def, entry.Def); err != nil {
			return IndexID{}, err
		}
	}

	newEntry := IndexEntry[T]{Def: def, Item: t}
	l.Exact[def.Signature] = append(l.Exact[def.Signature], newEntry)
	return IndexID{sig: def.Signature, idx: len(l.Exact[def.Signature]) - 1}, nil
}

func checkOverlap(left, right *Definition) error {
	if satisfiesContext(left.Context, right.Context) {
		return fmt.Errorf("Overlapping signatures found. All sources and contexts of this method\n    %s%s\n\nare contained in method\n    %s%s\n\nGoverter doesn't know which method to use when all contexts of the second method are available.\nRemove one of the methods to prevent this ambiguity.", left.ID, left.ArgDebug("        "), right.ID, right.ArgDebug("        "))
	}
	return nil
}

func (l *Index[T]) ByID(id IndexID) *T {
	return l.Exact[id.sig][id.idx].Item
}

func (l *Index[T]) Has(sig xtype.Signature) bool {
	_, ok := l.Exact[sig]
	return ok
}

func (l *Index[T]) Get(sig xtype.Signature, m map[string]*xtype.Type) (*T, error) {
	hits, ok := l.Exact[sig]
	if !ok {
		return nil, nil
	}

	for _, hit := range hits {
		if satisfiesContext(hit.Def.Context, m) {
			return hit.Item, nil
		}
	}

	return nil, satisfiedError(sig, m, hits)
}

func satisfiedError[T any](sig xtype.Signature, available map[string]*xtype.Type, hits []IndexEntry[T]) error {
	var hitStrings []string
	for _, hit := range hits {
		hitStrings = append(hitStrings, fmt.Sprintf("%s:\n    %s", hit.Def.ID, strings.Join(AvailableContextDebug(hit.Def.Context, available), "\n    ")))
	}
	return fmt.Errorf(`Found custom functions(s) converting %s to %s
but not all required context params are available in the current method.

%s`, sig.Source, sig.Target, strings.Join(hitStrings, "\n\n"))
}

func AvailableContextDebug(required, available map[string]*xtype.Type) []string {
	var lines []string

	use := xtype.UsageFromMap(available)
	for key := range required {
		_, ok := available[key]
		use.Used(key)

		prefix := "[available] "
		if !ok {
			prefix = "[missing]   "
		}
		lines = append(lines, prefix+key)
	}

	for _, x := range use.Unused() {
		lines = append(lines, "[unused]    "+x)
	}

	sort.Strings(lines)

	return lines
}

func satisfiesContext(required, m map[string]*xtype.Type) bool {
	for key := range required {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

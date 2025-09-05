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
	sig    xtype.Signature
	idx    int
	update bool
}

func NewIndex[T any]() *Index[T] {
	return &Index[T]{
		Exact: xtype.NewSignatureMap[[]IndexEntry[T]](),
	}
}

type Index[T any] struct {
	Exact  *xtype.SignatureMap[[]IndexEntry[T]]
	Update []*T
}

func (l *Index[T]) GetAll() []*T {
	items := []*T{}
	for _, exacts := range l.Exact.Values() {
		for _, exact := range exacts {
			items = append(items, exact.Item)
		}
	}
	return append(items, l.Update...)
}

func (l *Index[T]) RegisterOverrideOverlapping(t *T, def *Definition) {
	newEntry := IndexEntry[T]{Def: def, Item: t}
	entries, _ := l.Exact.At(def.Signature)
	for i, entry := range entries {
		if satisfiesContext(entry.Def.Context, def.Context) || satisfiesContext(def.Context, entry.Def.Context) {
			entries[i] = newEntry
			l.Exact.Set(def.Signature, entries)
			return
		}
	}

	entries = append(entries, newEntry)
	l.Exact.Set(def.Signature, entries)
}

func (l *Index[T]) RegisterUpdate(t *T, def *Definition) (IndexID, error) {
	l.Update = append(l.Update, t)
	return IndexID{update: true, idx: len(l.Update) - 1}, nil
}

func (l *Index[T]) Register(t *T, def *Definition) (IndexID, error) {
	entries, _ := l.Exact.At(def.Signature)
	for _, entry := range entries {
		if err := checkOverlap(entry.Def, def); err != nil {
			return IndexID{}, err
		}
		if err := checkOverlap(def, entry.Def); err != nil {
			return IndexID{}, err
		}
	}

	newEntry := IndexEntry[T]{Def: def, Item: t}
	entries = append(entries, newEntry)
	l.Exact.Set(def.Signature, entries)

	return IndexID{sig: def.Signature, idx: len(entries) - 1}, nil
}

func checkOverlap(left, right *Definition) error {
	if satisfiesContext(left.Context, right.Context) {
		return fmt.Errorf("Overlapping signatures found. All sources and contexts of this method\n    %s%s\n\nare contained in method\n    %s%s\n\nGoverter doesn't know which method to use when all contexts of the second method are available.\nRemove one of the methods to prevent this ambiguity.", left.ID, left.ArgDebug("        "), right.ID, right.ArgDebug("        "))
	}
	return nil
}

func (l *Index[T]) ByID(id IndexID) *T {
	if id.update {
		return l.Update[id.idx]
	}
	value, _ := l.Exact.At(id.sig)
	return value[id.idx].Item
}

func (l *Index[T]) Has(sig xtype.Signature) bool {
	_, ok := l.Exact.At(sig)
	return ok
}

func (l *Index[T]) Get(sig xtype.Signature, m map[string]*xtype.Type) (*T, error) {
	hits, ok := l.Exact.At(sig)
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

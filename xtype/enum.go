package xtype

import (
	"go/types"
	"sort"

	"github.com/jmattheis/goverter/enum"
)

type Enum struct {
	enum.Enum
	OK bool
}

func (t *Type) Enum(cfg *enum.Config) *Enum {
	if !t.Named {
		return disabled
	}

	if t.enum == nil {
		t.enum = loadEnum(t.NamedType, cfg)
	}
	return t.enum
}

func loadEnum(t *types.Named, cfg *enum.Config) *Enum {
	path := t.Obj().Pkg().Path()
	name := t.Obj().Name()

	if !cfg.Enabled || cfg.Excludes.Matches(path, name) {
		return disabled
	}

	e, ok := enum.Detect(complete(t, cfg))
	return &Enum{OK: ok, Enum: e}
}

// complete returns t as declared in its fully loaded package. A package that is only known
// through the export data of another package is incomplete: its scope holds the objects the
// other package references, so an enum's constants are missing in part or in full.
func complete(t *types.Named, cfg *enum.Config) *types.Named {
	pkg := t.Obj().Pkg()
	if pkg.Complete() || cfg.Package == nil {
		return t
	}
	full := cfg.Package(pkg.Path())
	if full == nil {
		return t
	}
	obj, ok := full.Scope().Lookup(t.Obj().Name()).(*types.TypeName)
	if !ok {
		return t
	}
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return t
	}
	return named
}

func (e Enum) SortedMembers() []string {
	var m []string
	for member := range e.Members {
		m = append(m, member)
	}
	sort.Strings(m)
	return m
}

var disabled = &Enum{OK: false}

package generator

import (
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

// loadPackages is used to load extend packages, with caching support.
func (g *generator) loadPackages(pkgPath string) ([]*packages.Package, error) {
	if pkgs, ok := g.pkgCache[pkgPath]; ok {
		return pkgs, nil
	}

	packagesCfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(packagesCfg, pkgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load package \"%s\", please add blank import for it", pkgPath)
	}

	if g.pkgCache == nil {
		g.pkgCache = make(map[string][]*packages.Package)
	}
	g.pkgCache[pkgPath] = pkgs
	return pkgs, nil
}

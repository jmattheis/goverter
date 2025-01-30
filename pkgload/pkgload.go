package pkgload

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jmattheis/goverter/config/parse"
	"github.com/jmattheis/goverter/method"
	"golang.org/x/tools/go/packages"
)

func New(workDir, buildTags string, paths []string) (*PackageLoader, error) {
	loader := &PackageLoader{
		lookupPkgPath: map[string]*packages.Package{},
		lookupAbsDir:  map[string]*packages.Package{},
		locals:        map[string]map[string]method.LocalOpts{},
	}
	_, err := loader.loadIntoCache(workDir, buildTags, paths)
	return loader, err
}

type PackageLoader struct {
	lookupPkgPath map[string]*packages.Package
	lookupAbsDir  map[string]*packages.Package
	locals        map[string]map[string]method.LocalOpts
}

func (g *PackageLoader) GetMatching(cwd, fullMethod string, opts *method.ParseOpts) ([]*method.Definition, error) {
	pkgName, name, err := ParseMethodString(cwd, fullMethod)
	if err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile(name)
	if err != nil {
		return nil, fmt.Errorf("could not parse name as regexp %q: %s", name, err)
	}

	if _, complete := pattern.LiteralPrefix(); complete {
		// not a regex
		m, err := g.getOneParsed(pkgName, name, opts)
		if err != nil {
			return nil, err
		}
		return []*method.Definition{m}, nil
	}

	// this is regexp, scan thru the package methods to find funcs that match the pattern
	var matches []*method.Definition

	pkg, err := g.getPkg(pkgName)
	if err != nil {
		return nil, err
	}

	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		loc := pattern.FindStringIndex(name)
		if len(loc) != 2 {
			continue
		}
		if loc[0] != 0 || loc[1] != len(name) {
			// we want full match only: e.g. CopyAbc.* won't match OtherCopyAbc
			continue
		}

		obj := scope.Lookup(name)
		m, err := method.Parse(obj, opts, g.localConfig(pkg, name))
		if err == nil {
			matches = append(matches, m)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf(`package %s does not have methods with names that match
the golang regexp pattern %q and a convert signature`, pkgName, name)
	}

	return matches, nil
}

func (g *PackageLoader) getPkg(pkgName string) (*packages.Package, error) {
	pkg := g.lookupPkgPath[pkgName]
	if pkg == nil {
		return nil, fmt.Errorf("failed to load package %q:\nmake sure it's a valid golang package", pkgName)
	}

	if len(pkg.Errors) > 0 {
		var lines []string
		for _, err := range pkg.Errors {
			lines = append(lines, err.Error())
		}

		return nil, fmt.Errorf("failed to load package %q:\n%s", pkgName, strings.Join(lines, "\n"))
	}
	return pkg, nil
}

func (g *PackageLoader) localConfig(pkg *packages.Package, name string) method.LocalOpts {
	fns, ok := g.locals[pkg.PkgPath]
	if !ok {
		fns = map[string]method.LocalOpts{}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					lines := parse.SettingLines(fn.Doc.Text())
					if len(lines) == 0 {
						continue
					}

					contexts := map[string]bool{}
					for _, line := range lines {
						if cmd, rest := parse.Command(line); cmd == "context" {
							if ctx, err := parse.String(rest); err == nil {
								contexts[ctx] = true
							}
						}
					}
					fns[fn.Name.Name] = method.LocalOpts{Context: contexts}
				}
			}
		}
		g.locals[pkg.PkgPath] = fns
	}
	fn, ok := fns[name]
	if !ok {
		return method.EmptyLocalOpts
	}
	return fn
}

func (g *PackageLoader) LoadPkgPathFromDir(dir string) (string, string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}
	if pkg, ok := g.lookupAbsDir[absDir]; ok {
		return pkg.Name, pkg.PkgPath, nil
	}
	// Skipping build tags as they're not used when only using packages.NeedName
	pkgs, err := load(absDir, "", packages.NeedName, []string{"."})
	if err != nil {
		return "", "", err
	}
	if len(pkgs) == 0 {
		return "", "", fmt.Errorf("no packages found in directory %s", absDir)
	}
	if len(pkgs) > 1 {
		return "", "", fmt.Errorf("too many packages found in the same directory %s", absDir)
	}
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return "", "", fmt.Errorf("got %d package errors; first error: %w",
			len(pkg.Errors),
			pkg.Errors[0])
	}
	g.lookupAbsDir[absDir] = pkg
	return pkg.Name, pkg.PkgPath, nil
}

func (g *PackageLoader) GetOneRaw(pkgName, name string) (*packages.Package, types.Object, error) {
	pkg, err := g.getPkg(pkgName)
	if err != nil {
		return pkg, nil, err
	}

	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return pkg, nil, fmt.Errorf("%q does not exist in package %q", name, pkgName)
	}
	return pkg, obj, nil
}

func (g *PackageLoader) GetOne(cwd, fullMethod string, opts *method.ParseOpts) (*method.Definition, error) {
	pkgName, name, err := ParseMethodString(cwd, fullMethod)
	if err != nil {
		return nil, err
	}
	return g.getOneParsed(pkgName, name, opts)
}

func (g *PackageLoader) getOneParsed(pkgName, name string, opts *method.ParseOpts) (*method.Definition, error) {
	pkg, obj, err := g.GetOneRaw(pkgName, name)
	if err != nil {
		return nil, err
	}

	def, err := method.Parse(obj, opts, g.localConfig(pkg, name))
	if err != nil {
		return nil, err
	}
	return def, nil
}

// loadIntoCache is used to load extend packages, with caching support.
func (g *PackageLoader) loadIntoCache(workDir, buildTags string, paths []string) ([]*packages.Package, error) {
	mode := packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax
	pkgs, err := load(workDir, buildTags, mode, paths)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		g.lookupPkgPath[pkg.PkgPath] = pkg
		if absDir := packageAbsDir(pkg); absDir != "" {
			g.lookupAbsDir[absDir] = pkg
		}
	}
	return pkgs, nil
}

func packageAbsDir(pkg *packages.Package) string {
	var dir string
	pkg.Fset.Iterate(func(f *token.File) bool {
		if filepath.IsAbs(f.Name()) {
			dir = filepath.Dir(f.Name())
			return false // stop iterating
		}
		return true // continue
	})
	return dir
}

// load is used to load extend packages and referenced output directory packages
func load(workDir, buildTags string, mode packages.LoadMode, paths []string) ([]*packages.Package, error) {
	packagesCfg := &packages.Config{
		Mode: mode,
		Dir:  workDir,
	}
	if buildTags != "" {
		packagesCfg.BuildFlags = append(packagesCfg.BuildFlags, "-tags", buildTags)
	}
	pkgs, err := packages.Load(packagesCfg, paths...)
	if err != nil {
		// This happens rare, and only if somebody uses advanced package pattern query in a wrong way.
		// The cause (err) usually has enough details to troubleshoot this issue.
		return nil, fmt.Errorf("failed to load packages %s:\n%s", paths, err)
	}
	return pkgs, nil
}

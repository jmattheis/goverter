package config

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"

	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/pkgload"
)

// ParseExtendOptions holds extend method options.
type ParseExtendOptions struct {
	// PkgPath where the extend methods are located. If it is empty, the package is same as the
	// ConverterInterface package and ConverterScope should be used for the lookup.
	PkgPath string
	// Scope of the ConverterInterface.
	ConverterScope *types.Scope
	// ConverterInterface to use - can be nil if its use is not allowed.
	ConverterInterface types.Type
	// NamePattern is the regexp pattern to search for within the PkgPath above or
	// (if PkgPath is empty) within the Scope.
	NamePattern *regexp.Regexp
}

// parseExtend prepares a list of extend methods for use.
func parseExtend(loader *pkgload.PackageLoader, c *Converter, methods []string) ([]*method.Definition, error) {
	var defs []*method.Definition
	for _, methodName := range methods {
		pkgPath, namePattern, err := splitCustomMethod(methodName)
		if err != nil {
			return nil, err
		}

		pattern, err := regexp.Compile(namePattern)
		if err != nil {
			return nil, fmt.Errorf("could not parse name as regexp %q: %s", namePattern, err)
		}

		opts := &ParseExtendOptions{
			PkgPath:            pkgPath,
			NamePattern:        pattern,
			ConverterInterface: c.Type,
			ConverterScope:     c.Scope,
		}

		methodDefs, err := parseExtendPackage(loader, opts)
		defs = append(defs, methodDefs...)
		if err != nil {
			return nil, err
		}
	}
	return defs, nil
}

// parseExtendPackage parses the goverter:extend inputs with or without packages (local or external).
//
// extend statement can be one of the following:
// 1) local scope with a name: "ConvertAToB", it is also equivalent to ":ConvertAToB"
// 2) package with a name: "github.com/google/uuid:FromBytes"
// 3) either (1) or (2) with the above with a regexp pattern instead of a name
//
// To scan the whole package for candidate methods, use "package/path:.*".
// Note: if regexp pattern is used, only the methods matching the conversion signature can be used.
// Those are methods that have exactly one input (to convert from) and either one output (to covert to)
// or two outputs: type to convert to and an error object.
func parseExtendPackage(loader *pkgload.PackageLoader, opts *ParseExtendOptions) ([]*method.Definition, error) {
	if opts.PkgPath == "" {
		// search in the converter's scope
		defs, err := searchExtendsInScope(opts.ConverterScope, opts)
		if err == nil && len(defs) == 0 {
			// no failure, but also nothing found (this can happen if pattern is used yet no matches found)
			err = fmt.Errorf("local package does not have methods with names that match "+
				"the golang regexp pattern %q and a convert signature", opts.NamePattern)
		}
		return defs, err
	}

	return searchExtendsInPackages(loader, opts)
}

// searchExtendsInPackages searches for extend conversion methods that match an input regexp pattern
// within a given package path.
// Note: if this method finds no candidates, it will report an error. Two reasons for that:
// scanning packages takes time and it is very likely a human error.
func searchExtendsInPackages(loader *pkgload.PackageLoader, opts *ParseExtendOptions) ([]*method.Definition, error) {
	// load a package by its path, loadPackages uses cache
	pkgs, err := loader.Load(opts.PkgPath)
	if err != nil {
		return nil, err
	}

	var defs []*method.Definition
	for _, pkg := range pkgs {
		// search in the scope of each package, first package is going to be the root one
		pkgDefs, pkgErr := searchExtendsInScope(pkg.Types.Scope(), opts)
		defs = append(defs, pkgDefs...)
		if pkgErr != nil && err == nil {
			err = pkgErr
		}
	}

	if len(defs) == 0 {
		if err == nil {
			return nil, fmt.Errorf(`package %s does not have methods with names that match
the golang regexp pattern %q and a convert signature`, opts.PkgPath, opts.NamePattern.String())
		}
		return nil, err
	}

	return defs, nil
}

// searchExtendsInScope searches the given package scope (either local or external) for
// the conversion method candidates. See parseExtendPackage for more details.
// If the input scope is not local, always pass converterInterface as a nil.
func searchExtendsInScope(scope *types.Scope, opts *ParseExtendOptions) ([]*method.Definition, error) {
	if prefix, complete := opts.NamePattern.LiteralPrefix(); complete {
		// this is not a regexp, use regular lookup and report error as is
		// we expect only one function to match
		def, err := parseExtendScopeMethod(scope, prefix, opts)
		if err != nil {
			return nil, err
		}
		return []*method.Definition{def}, nil
	}

	// this is regexp, scan thru the package methods to find funcs that match the pattern
	var defs []*method.Definition
	for _, name := range scope.Names() {
		loc := opts.NamePattern.FindStringIndex(name)
		if len(loc) != 2 {
			continue
		}
		if loc[0] != 0 || loc[1] != len(name) {
			// we want full match only: e.g. CopyAbc.* won't match OtherCopyAbc
			continue
		}

		def, err := parseExtendScopeMethod(scope, name, opts)
		if err == nil {
			defs = append(defs, def)
		}
	}
	return defs, nil
}

func splitCustomMethod(fullMethod string) (path, name string, err error) {
	parts := strings.SplitN(fullMethod, ":", 2)
	switch len(parts) {
	case 0:
		return "", "", fmt.Errorf("invalid custom method: %s", fullMethod)
	case 1:
		name = parts[0]
	case 2:
		path = parts[0]
		name = parts[1]
		if path == "" {
			// example: goverter:extend :MyLocalConvert
			// the purpose of the ':' in this case is confusing, do not allow such case
			return "", "", fmt.Errorf(`package path must not be empty in the custom method "%s".
See https://goverter.jmattheis.de/#/config/extend`, fullMethod)
		}
	}

	if name == "" {
		return "", "", fmt.Errorf(`method name pattern is required in the custom method "%s".
See https://goverter.jmattheis.de/#/config/extend`, fullMethod)
	}
	return
}

// parseExtend prepares an extend conversion method using its name and a scope to search.
func parseExtendScopeMethod(scope *types.Scope, methodName string, opts *ParseExtendOptions) (*method.Definition, error) {
	obj := scope.Lookup(methodName)
	if obj == nil {
		return nil, fmt.Errorf("%s does not exist in scope", methodName)
	}
	return method.Parse(&method.ParseOpts{
		ErrorPrefix: "error parsing type",
		Obj:         obj,
		Converter:   opts.ConverterInterface,
		EmptySource: false,
	})
}

// parseExtend prepares an extend conversion method using its name and a scope to search.
func parseMapExtend(loader *pkgload.PackageLoader, c *Converter, fullMethod string) (*method.Definition, error) {
	pkgPath, name, err := splitCustomMethod(fullMethod)
	if err != nil {
		return nil, err
	}

	useScope := c.Scope

	if pkgPath != "" {
		pkgs, err := loader.Load(pkgPath)
		if err != nil {
			return nil, err
		}
		if len(pkgs) != 1 {
			return nil, fmt.Errorf("'%s' package path matches multiple packages, it must match exactly one", fullMethod)
		}
		useScope = pkgs[0].Types.Scope()
	}

	obj := useScope.Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("%s does not exist in scope", name)
	}

	return method.Parse(&method.ParseOpts{
		ErrorPrefix: "error parsing type",
		Obj:         obj,
		Converter:   c.Type,
		EmptySource: true,
	})
}

package config

import (
	"path/filepath"
	"strings"

	"github.com/jmattheis/goverter/config/parse"
	"github.com/jmattheis/goverter/pkgload"
)

func getPackages(raw *Raw) []string {
	lookup := map[string]struct{}{}
	for _, c := range raw.Converters {
		lookup[c.PackagePath] = struct{}{}

		// the default output:file is in ./generated and may not be configured in Raw.
		// This preemptively loads this package, in case it already exists.
		lookup[filepath.Join(c.PackagePath, "generated")] = struct{}{}

		registerConverterLines(lookup, c.PackagePath, c.Converter)
		registerConverterLines(lookup, c.PackagePath, raw.Global)
		for _, m := range c.Methods {
			registerMethodLines(lookup, c.PackagePath, m)
		}
	}

	var pkgs []string
	for pkg := range lookup {
		pkgs = append(pkgs, "pattern="+pkg)
	}

	return pkgs
}

func registerConverterLines(lookup map[string]struct{}, cwd string, lines RawLines) {
	for _, line := range lines.Lines {
		cmd, rest := parse.Command(line)
		switch cmd {
		case configExtend:
			for _, fullMethod := range strings.Fields(rest) {
				registerFullMethod(lookup, cwd, fullMethod)
			}
		case configOutputFile:
			if file, err := parse.String(rest); err == nil {
				lookup[filepath.Dir(pkgload.ResolveRelativePath(cwd, file))] = struct{}{}
			}
		}
	}
}

func registerMethodLines(lookup map[string]struct{}, cwd string, lines RawLines) {
	for _, line := range lines.Lines {
		cmd, rest := parse.Command(line)
		switch cmd {
		case configMap:
			if _, _, custom, err := parseMethodMap(rest); err == nil && custom != "" {
				registerFullMethod(lookup, cwd, custom)
			}
		case configDefault:
			registerFullMethod(lookup, cwd, rest)
		}
	}
}

func registerFullMethod(lookup map[string]struct{}, cwd, fullMethod string) {
	pkg, _, err := pkgload.ParseMethodString(cwd, fullMethod)
	if err == nil {
		lookup[pkg] = struct{}{}
	}
}

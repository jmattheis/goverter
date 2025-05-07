package config

import (
	"path/filepath"
	"strings"

	"github.com/jmattheis/goverter/config/parse"
	"github.com/jmattheis/goverter/pkgload"
)

func resolvePackage(sourceFileName, sourcePackage, targetFile string) (string, error) {
	relativeFile := targetFile
	if filepath.IsAbs(targetFile) {
		var err error
		relativeFile, err = filepath.Rel(filepath.Dir(sourceFileName), targetFile)
		if err != nil {
			return "", err
		}
	}

	return filepath.ToSlash(filepath.Dir(filepath.Join(sourcePackage, relativeFile))), nil
}

func getPackages(raw *Raw) []string {
	lookup := map[string]struct{}{}
	for _, c := range raw.Converters {
		lookup[c.PackagePath] = struct{}{}

		// the default output:file is in ./generated and is not configured in Raw.
		// This preemptively loads this package, in case it already exists.
		lookup[filepath.Join(c.PackagePath, "generated")] = struct{}{}

		registerConverterLines(lookup, raw.WorkDir, c.FileName, c.PackagePath, c.Converter)
		registerConverterLines(lookup, raw.WorkDir, c.FileName, c.PackagePath, raw.Global)
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

func registerConverterLines(lookup map[string]struct{}, cwd, filename, sourcePackage string, lines RawLines) {
	for _, line := range lines.Lines {
		cmd, rest := parse.Command(line)
		switch cmd {
		case configExtend:
			for _, fullMethod := range strings.Fields(rest) {
				registerFullMethod(lookup, sourcePackage, fullMethod)
			}
		case configOutputFile:
			file, err := parse.File(cwd, rest)
			if err != nil {
				continue
			}
			targetPackage, err := resolvePackage(filename, sourcePackage, file)
			if err != nil {
				continue
			}
			lookup[targetPackage] = struct{}{}
		}
	}
}

func registerMethodLines(lookup map[string]struct{}, sourcePackage string, lines RawLines) {
	for _, line := range lines.Lines {
		cmd, rest := parse.Command(line)
		switch cmd {
		case configMap:
			if _, _, custom, err := parseMethodMap(rest); err == nil && custom != "" {
				registerFullMethod(lookup, sourcePackage, custom)
			}
		case configDefault:
			registerFullMethod(lookup, sourcePackage, rest)
		}
	}
}

func registerFullMethod(lookup map[string]struct{}, sourcePackage, fullMethod string) {
	pkg, _, err := pkgload.ParseMethodString(sourcePackage, fullMethod)
	if err == nil {
		lookup[pkg] = struct{}{}
	}
}

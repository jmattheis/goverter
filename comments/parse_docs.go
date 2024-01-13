package comments

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/jmattheis/goverter/config"
	"golang.org/x/tools/go/packages"
)

const (
	prefix          = "goverter"
	delimter        = ":"
	converterMarker = prefix + delimter + "converter"
	variablesMarker = prefix + delimter + "variables"
)

// ParseDocsConfig provides input to the ParseDocs method below.
type ParseDocsConfig struct {
	// PackagePatterns are golang package patterns to scan, required.
	PackagePattern []string
	// WorkingDir is a directory to invoke the tool on. If omitted, current directory is used.
	WorkingDir string
	BuildTags  string
}

// Converter defines a converter that was marked with converterMarker.
type Converter struct {
	Name     string
	Comments config.RawLines
	Methods  map[string]config.RawLines
	Scope    *types.Scope
}

// ParseDocs parses the docs for the given pattern.
func ParseDocs(c ParseDocsConfig) ([]config.RawConverter, error) {
	loadCfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Dir:  c.WorkingDir,
	}
	if c.BuildTags != "" {
		loadCfg.BuildFlags = append(loadCfg.BuildFlags, "-tags", c.BuildTags)
	}
	pkgs, err := packages.Load(loadCfg, c.PackagePattern...)
	if err != nil {
		return nil, err
	}
	rawConverters := []config.RawConverter{}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf(`could not load package %s

%s

Goverter cannot generate converters when there are compile errors because it
requires the type information from the compiled sources.`, pkg.PkgPath, pkg.Errors[0])
		}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					converters, err := parseGenDecl(pkg.Fset, pkg.Types, genDecl)
					if err != nil {
						location := pkg.Fset.Position(genDecl.Pos()).String()
						return rawConverters, fmt.Errorf("%s: %s", location, err)
					}
					rawConverters = append(rawConverters, converters...)
				}
			}
		}
	}
	return rawConverters, nil
}

func parseFunctions(fset *token.FileSet, pkg *types.Package, decl *ast.GenDecl, comments string) ([]config.RawConverter, error) {
	if decl.Tok != token.VAR {
		return nil, fmt.Errorf("%s must be defined on %q-block but was %q", converterMarker, token.VAR, decl.Tok.String())
	}

	location := fset.Position(decl.Pos())
	converterLines := parseRawLines(fileWithLine(location), comments)

	result := map[string]config.RawLines{}
	for _, spec := range decl.Specs {
		value, ok := spec.(*ast.ValueSpec)
		if !ok {
			return nil, fmt.Errorf("expected value spec but got %#v", spec)
		}
		if len(value.Names) != 1 {
			return nil, fmt.Errorf("must have one name")
		}
		name := value.Names[0].Name

		location := fileWithLine(fset.Position(value.Pos()))
		result[name] = parseRawLines(location, value.Doc.Text())
	}

	converter := config.RawConverter{
		FileName:    location.Filename,
		Converter:   converterLines,
		Methods:     result,
		PackageName: pkg.Name(),
		PackagePath: pkg.Path(),
	}
	return []config.RawConverter{converter}, nil
}

func parseGenDecl(fset *token.FileSet, pkg *types.Package, decl *ast.GenDecl) ([]config.RawConverter, error) {
	declDocs := decl.Doc.Text()

	if strings.Contains(declDocs, variablesMarker) {
		return parseFunctions(fset, pkg, decl, declDocs)
	}

	if strings.Contains(declDocs, converterMarker) {
		if decl.Tok != token.TYPE {
			return nil, fmt.Errorf("%s must be defined on %q-block but was %q", converterMarker, token.TYPE, decl.Tok.String())
		}

		if len(decl.Specs) != 1 {
			return nil, fmt.Errorf("found %s on type but it has multiple interfaces inside", converterMarker)
		}
		typeSpec, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return nil, fmt.Errorf("%s may only be applied to type declarations ", converterMarker)
		}
		c, err := parseInterface(fset, pkg, typeSpec, declDocs)
		if err != nil {
			return nil, err
		}
		return []config.RawConverter{c}, nil
	}

	var converters []config.RawConverter

	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok && strings.Contains(typeSpec.Doc.Text(), converterMarker) {
			c, err := parseInterface(fset, pkg, typeSpec, typeSpec.Doc.Text())
			if err != nil {
				return nil, err
			}
			converters = append(converters, c)
		}
	}

	return converters, nil
}

func parseInterface(fset *token.FileSet, pkg *types.Package, typeSpec *ast.TypeSpec, declDocs string) (config.RawConverter, error) {
	astInterface, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return config.RawConverter{}, fmt.Errorf("%s may only be applied to type interface declarations ", converterMarker)
	}
	typeName := typeSpec.Name.String()

	location := fset.Position(typeSpec.Pos())
	converterLines := parseRawLines(fileWithLine(location), declDocs)
	methods, err := parseInterfaceMethods(fset, astInterface)
	if err != nil {
		return config.RawConverter{}, fmt.Errorf("type %s: %s", typeName, err)
	}
	converter := config.RawConverter{
		InterfaceName: typeName,
		FileName:      location.Filename,
		Converter:     converterLines,
		Methods:       methods,
		PackageName:   pkg.Name(),
		PackagePath:   pkg.Path(),
	}
	return converter, nil
}

func parseInterfaceMethods(location *token.FileSet, inter *ast.InterfaceType) (map[string]config.RawLines, error) {
	result := map[string]config.RawLines{}
	for _, method := range inter.Methods.List {
		if len(method.Names) != 1 {
			return result, fmt.Errorf("method must have one name")
		}
		name := method.Names[0].String()

		location := location.Position(method.Pos())
		result[name] = parseRawLines(fileWithLine(location), method.Doc.Text())
	}
	return result, nil
}

func parseRawLines(location, comment string) config.RawLines {
	scanner := bufio.NewScanner(strings.NewReader(comment))
	raw := config.RawLines{Location: location}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, prefix+delimter) {
			line := strings.TrimPrefix(line, prefix+delimter)
			raw.Lines = append(raw.Lines, line)
		}
	}
	return raw
}

func fileWithLine(p token.Position) string {
	return fmt.Sprintf("%s:%d", p.Filename, p.Line)
}

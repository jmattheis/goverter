package comments

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	prefix          = "goverter"
	delimter        = ":"
	converterMarker = prefix + delimter + "converter"
)

// MethodMapping a mapping between method name and method.
type MethodMapping map[string]Method

// ParseDocsConfig provides input to the ParseDocs method below.
type ParseDocsConfig struct {
	// PackagePattern is a golang package pattern to scan, required.
	PackagePattern string
	// WorkingDir is a directory to invoke the tool on. If omitted, current directory is used.
	WorkingDir string
}

// Converter defines a converter that was marked with converterMarker.
type Converter struct {
	Name    string
	Config  ConverterConfig
	Methods MethodMapping
	Scope   *types.Scope
}

// ConverterConfig contains settings that can be set via comments.
type ConverterConfig struct {
	Name          string
	ExtendMethods []string
	WrapErrors    bool
}

// Method contains settings that can be set via comments.
type Method struct {
	IgnoredFields   map[string]struct{}
	NameMapping     map[string]string
	MatchIgnoreCase bool
	WrapErrors      bool
	// target to source
	IdentityMapping map[string]struct{}
	// mapping function to source
	ExtendMapping map[string]string
}

// ParseDocs parses the docs for the given pattern.
func ParseDocs(config ParseDocsConfig) ([]Converter, error) {
	loadCfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedCompiledGoFiles | packages.NeedTypes |
			packages.NeedModule | packages.NeedFiles | packages.NeedName | packages.NeedImports,
		Dir: config.WorkingDir,
	}
	pkgs, err := packages.Load(loadCfg, config.PackagePattern)
	if err != nil {
		return nil, err
	}
	mapping := []Converter{}
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
					converters, err := parseGenDecl(pkg.Types.Scope(), genDecl)
					if err != nil {
						return mapping, fmt.Errorf("%s: %s", pkg.Fset.Position(genDecl.Pos()), err)
					}
					mapping = append(mapping, converters...)
				}
			}
		}
	}
	sort.Slice(mapping, func(i, j int) bool {
		return mapping[i].Config.Name < mapping[j].Config.Name
	})
	return mapping, nil
}

func parseGenDecl(scope *types.Scope, decl *ast.GenDecl) ([]Converter, error) {
	declDocs := decl.Doc.Text()

	if strings.Contains(declDocs, converterMarker) {
		if len(decl.Specs) != 1 {
			return nil, fmt.Errorf("found %s on type but it has multiple interfaces inside", converterMarker)
		}
		typeSpec, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return nil, fmt.Errorf("%s may only be applied to type declarations ", converterMarker)
		}
		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return nil, fmt.Errorf("%s may only be applied to type interface declarations ", converterMarker)
		}
		typeName := typeSpec.Name.String()
		config, err := parseConverterComment(declDocs, ConverterConfig{Name: typeName + "Impl"})
		if err != nil {
			return nil, fmt.Errorf("type %s: %s", typeName, err)
		}
		methods, err := parseInterface(interfaceType)
		if err != nil {
			return nil, fmt.Errorf("type %s: %s", typeName, err)
		}
		converter := Converter{
			Name:    typeName,
			Methods: methods,
			Config:  config,
			Scope:   scope,
		}
		return []Converter{converter}, nil
	}

	var converters []Converter

	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok && strings.Contains(typeSpec.Doc.Text(), converterMarker) {
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return nil, fmt.Errorf("%s may only be applied to type interface declarations ", converterMarker)
			}
			typeName := typeSpec.Name.String()
			config, err := parseConverterComment(typeSpec.Doc.Text(), ConverterConfig{Name: typeName + "Impl"})
			if err != nil {
				return nil, fmt.Errorf("type %s: %s", typeName, err)
			}
			methods, err := parseInterface(interfaceType)
			if err != nil {
				return nil, fmt.Errorf("type %s: %s", typeName, err)
			}
			converters = append(converters, Converter{
				Name:    typeName,
				Methods: methods,
				Config:  config,
				Scope:   scope,
			})
		}
	}

	return converters, nil
}

func parseInterface(inter *ast.InterfaceType) (MethodMapping, error) {
	result := MethodMapping{}
	for _, method := range inter.Methods.List {
		if len(method.Names) != 1 {
			return result, fmt.Errorf("method must have one name")
		}
		name := method.Names[0].String()

		parsed, err := parseMethodComment(method.Doc.Text())
		if err != nil {
			return result, fmt.Errorf("parsing method %s: %s", name, err)
		}

		result[name] = parsed
	}
	return result, nil
}

func parseConverterComment(comment string, config ConverterConfig) (ConverterConfig, error) {
	scanner := bufio.NewScanner(strings.NewReader(comment))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, prefix+delimter) {
			cmd := strings.TrimPrefix(line, prefix+delimter)
			if cmd == "" {
				return config, fmt.Errorf("unknown %s comment: %s", prefix, line)
			}
			fields := strings.Fields(cmd)
			switch fields[0] {
			case "converter":
				// only a marker interface
				continue
			case "name":
				if len(fields) != 2 {
					return config, fmt.Errorf("invalid %s:name must have one parameter", prefix)
				}
				config.Name = fields[1]
				continue
			case "extend":
				config.ExtendMethods = append(config.ExtendMethods, fields[1:]...)
				continue
			case "wrapErrors":
				if len(fields) != 1 {
					return config, fmt.Errorf("invalid %s:wrapErrors, parameters not supported", prefix)
				}
				config.WrapErrors = true
				continue
			}
			return config, fmt.Errorf("unknown %s comment: %s", prefix, line)
		}
	}
	return config, nil
}

func parseMethodComment(comment string) (Method, error) {
	scanner := bufio.NewScanner(strings.NewReader(comment))
	m := Method{
		NameMapping:     map[string]string{},
		IgnoredFields:   map[string]struct{}{},
		IdentityMapping: map[string]struct{}{},
		ExtendMapping:   map[string]string{},
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, prefix+delimter) {
			cmd := strings.TrimPrefix(line, prefix+delimter)
			if cmd == "" {
				return m, fmt.Errorf("unknown %s comment: %s", prefix, line)
			}
			fields := strings.Fields(cmd)
			switch fields[0] {
			case "map":
				if len(fields) != 3 {
					return m, fmt.Errorf("invalid %s:map must have two parameter", prefix)
				}
				m.NameMapping[fields[2]] = fields[1]
				continue
			case "mapIdentity":
				for _, f := range fields[1:] {
					m.IdentityMapping[f] = struct{}{}
				}
				continue
			case "ignore":
				for _, f := range fields[1:] {
					m.IgnoredFields[f] = struct{}{}
				}
				continue
			case "mapExtend":
				if len(fields) != 3 {
					return m, fmt.Errorf("invalid %s:mapExtend must have two parameter", prefix)
				}
				m.ExtendMapping[fields[1]] = fields[2]
				continue
			case "matchIgnoreCase":
				if len(fields) != 1 {
					return m, fmt.Errorf("invalid %s:matchIgnoreCase, parameters not supported", prefix)
				}
				m.MatchIgnoreCase = true
				continue
			case "wrapErrors":
				if len(fields) != 1 {
					return m, fmt.Errorf("invalid %s:wrapErrors, parameters not supported", prefix)
				}
				m.WrapErrors = true
				continue
			}
			return m, fmt.Errorf("unknown %s comment: %s", prefix, line)
		}
	}
	return m, nil
}

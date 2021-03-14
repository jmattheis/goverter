package comments

import (
	"bufio"
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	prefix          = "genconv"
	delimter        = ":"
	converterMarker = prefix + delimter + "converter"
)

type MethodMapping map[string]Method

type Converter struct {
	Name    string
	Config  ConverterConfig
	Methods MethodMapping
}

type ConverterConfig struct {
	Name          string
	ExtendMethods []string
}

type Method struct {
	IgnoredFields map[string]struct{}
	NameMapping   map[string]string
	// target to source
}

func ParseDocs(pattern string) ([]Converter, error) {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.LoadAllSyntax}, pattern)
	if err != nil {
		return nil, err
	}
	mapping := []Converter{}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					var err error
					if mapping, err = parseGenDecl(mapping, genDecl); err != nil {
						return mapping, fmt.Errorf("%s: %s", pkg.Fset.Position(genDecl.Pos()), err)
					}
				}
			}
		}
	}
	sort.Slice(mapping, func(i, j int) bool {
		return mapping[i].Config.Name < mapping[j].Config.Name
	})
	return mapping, nil
}

func parseGenDecl(mapping []Converter, decl *ast.GenDecl) ([]Converter, error) {
	declDocs := decl.Doc.Text()

	if strings.Contains(declDocs, converterMarker) {
		if len(decl.Specs) != 1 {
			return mapping, fmt.Errorf("found %s on type but it has multiple interfaces inside", converterMarker)
		}
		typeSpec, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return mapping, fmt.Errorf("%s may only be applied to type declarations ", converterMarker)
		}
		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return mapping, fmt.Errorf("%s may only be applied to type interface declarations ", converterMarker)
		}
		typeName := typeSpec.Name.String()
		config, err := parseConverterComment(declDocs, ConverterConfig{Name: typeName + "Impl"})
		if err != nil {
			return mapping, fmt.Errorf("type %s: %s", typeName, err)
		}
		methods, err := parseInterface(interfaceType)
		if err != nil {
			return mapping, fmt.Errorf("type %s: %s", typeName, err)
		}
		mapping = append(mapping, Converter{
			Name:    typeName,
			Methods: methods,
			Config:  config,
		})
		return mapping, nil
	}

	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok && strings.Contains(typeSpec.Doc.Text(), converterMarker) {
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return mapping, fmt.Errorf("%s may only be applied to type interface declarations ", converterMarker)
			}
			typeName := typeSpec.Name.String()
			config, err := parseConverterComment(typeSpec.Doc.Text(), ConverterConfig{Name: typeName + "Impl"})
			if err != nil {
				return mapping, fmt.Errorf("type %s: %s", typeName, err)
			}
			methods, err := parseInterface(interfaceType)
			if err != nil {
				return mapping, fmt.Errorf("type %s: %s", typeName, err)
			}
			mapping = append(mapping, Converter{
				Name:    typeName,
				Methods: methods,
				Config:  config,
			})
		}
	}

	return mapping, nil
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
			}
			return config, fmt.Errorf("unknown %s comment: %s", prefix, line)
		}
	}
	return config, nil
}

func parseMethodComment(comment string) (Method, error) {
	scanner := bufio.NewScanner(strings.NewReader(comment))
	m := Method{
		NameMapping:   map[string]string{},
		IgnoredFields: map[string]struct{}{},
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
			case "ignore":
				if len(fields) != 2 {
					return m, fmt.Errorf("invalid %s:ignore must have two parameter", prefix)
				}
				m.IgnoredFields[fields[1]] = struct{}{}
				continue
			}
			return m, fmt.Errorf("unknown %s comment: %s", prefix, line)
		}
	}
	return m, nil
}

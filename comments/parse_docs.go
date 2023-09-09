package comments

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"github.com/jmattheis/goverter/builder"
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
	Flags         builder.ConversionFlags
}

// Method contains settings that can be set via comments.
type Method struct {
	Flags   builder.ConversionFlags
	AutoMap []string
	Fields  map[string]*FieldMapping
}

func (m *Method) Field(targetName string) *FieldMapping {
	target, ok := m.Fields[targetName]
	if !ok {
		target = &FieldMapping{}
		m.Fields[targetName] = target
	}
	return target
}

type FieldMapping struct {
	Source   string
	Function string
	Ignore   bool
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
		config, err := parseConverterComment(declDocs, ConverterConfig{
			Name:  typeName + "Impl",
			Flags: builder.ConversionFlags{},
		})
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
			config, err := parseConverterComment(typeSpec.Doc.Text(), ConverterConfig{Name: typeName + "Impl", Flags: builder.ConversionFlags{}})
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
				config.Flags.Set(builder.FlagWrapErrors)
				continue
			case "ignoreUnexported":
				config.Flags.Set(builder.FlagIgnoreUnexported)
				continue
			case "matchIgnoreCase":
				config.Flags.Set(builder.FlagMatchIgnoreCase)
				continue
			case "ignoreMissing":
				config.Flags.Set(builder.FlagIgnoreMissing)
				continue
			case "skipCopySameType":
				config.Flags.Set(builder.FlagSkipCopySameType)
				continue
			case "useZeroValueOnPointerInconsistency":
				config.Flags.Set(builder.FlagZeroValueOnPtrInconsistency)
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
		Fields: map[string]*FieldMapping{},
		Flags:  builder.ConversionFlags{},
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, prefix+delimter) {
			cmd := strings.TrimPrefix(line, prefix+delimter)
			if cmd == "" {
				return m, fmt.Errorf("unknown %s comment: %s", prefix, line)
			}
			parts := strings.SplitN(cmd, " ", 2)
			key := parts[0]

			remaining := ""
			if len(parts) == 2 {
				remaining = parts[1]
			}
			switch key {
			case "map":
				parts := strings.SplitN(remaining, "|", 2)
				fields := strings.Fields(parts[0])
				custom := ""
				if len(parts) == 2 {
					custom = strings.TrimSpace(parts[1])
				}

				switch len(fields) {
				case 1:
					f := m.Field(fields[0])
					f.Function = custom
				case 2:
					f := m.Field(fields[1])
					f.Source = fields[0]
					f.Function = custom
				case 0:
					return m, fmt.Errorf("invalid %s:map missing target field", prefix)
				default:
					return m, fmt.Errorf("invalid %s:map too many fields", prefix)
				}
				continue
			case "mapIdentity":
				fields := strings.Fields(remaining)
				for _, f := range fields {
					m.Field(f).Source = "."
				}
				continue
			case "ignore":
				fields := strings.Fields(remaining)
				for _, f := range fields {
					m.Field(f).Ignore = true
				}
				continue
			case "autoMap":
				m.AutoMap = append(m.AutoMap, strings.TrimSpace(remaining))
				continue
			case "mapExtend":
				fields := strings.Fields(remaining)
				if len(fields) != 2 {
					return m, fmt.Errorf("invalid %s:mapExtend must have two parameter", prefix)
				}
				f := m.Field(fields[0])
				f.Function = fields[1]
				f.Source = "."
				continue
			case "matchIgnoreCase":
				if strings.TrimSpace(remaining) != "" {
					return m, fmt.Errorf("invalid %s:matchIgnoreCase, parameters not supported", prefix)
				}
				m.Flags.Set(builder.FlagMatchIgnoreCase)
				continue
			case "ignoreMissing":
				m.Flags.Set(builder.FlagIgnoreMissing)
				continue
			case "ignoreUnexported":
				m.Flags.Set(builder.FlagIgnoreUnexported)
				continue
			case "useZeroValueOnPointerInconsistency":
				m.Flags.Set(builder.FlagZeroValueOnPtrInconsistency)
				continue
			case "skipCopySameType":
				m.Flags.Set(builder.FlagSkipCopySameType)
				continue
			case "wrapErrors":
				if strings.TrimSpace(remaining) != "" {
					return m, fmt.Errorf("invalid %s:wrapErrors, parameters not supported", prefix)
				}
				m.Flags.Set(builder.FlagWrapErrors)
				continue
			}
			return m, fmt.Errorf("unknown %s comment: %s", prefix, line)
		}
	}
	return m, nil
}

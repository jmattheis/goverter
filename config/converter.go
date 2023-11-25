package config

import (
	"go/types"
	"strings"

	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/pkgload"
)

var DefaultConfig = ConverterConfig{
	OutputFile:        "./generated/generated.go",
	OutputPackageName: "generated",
}

type Converter struct {
	ConverterConfig
	FileSource string
	Type       types.Type
	Interface  *types.Interface
	Scope      *types.Scope
	Methods    map[string]*Method
}

type ConverterConfig struct {
	Common
	Name              string
	OutputFile        string
	OutputPackagePath string
	OutputPackageName string
	Extend            []*method.Definition
	Comments          []string
}

func (conf *ConverterConfig) PackageID() string {
	if conf.OutputPackageName == "" {
		return conf.OutputPackagePath
	}
	return conf.OutputPackagePath + ":" + conf.OutputPackageName
}

func parseGlobal(loader *pkgload.PackageLoader, global RawLines) (*ConverterConfig, error) {
	c := Converter{ConverterConfig: DefaultConfig}
	err := parseConverterLines(&c, "global", loader, global)
	return &c.ConverterConfig, err
}

func parseConverter(loader *pkgload.PackageLoader, rawConverter *RawConverter, global ConverterConfig) (*Converter, error) {
	namedType := rawConverter.Scope.Lookup(rawConverter.InterfaceName).Type()
	interfaceType := namedType.Underlying().(*types.Interface)

	c := &Converter{
		ConverterConfig: global,
		Type:            namedType,
		Interface:       interfaceType,
		Scope:           rawConverter.Scope,
		FileSource:      rawConverter.FileSource,
		Methods:         map[string]*Method{},
	}
	if c.Name == "" {
		c.Name = rawConverter.InterfaceName + "Impl"
	}

	if err := parseConverterLines(c, c.Type.String(), loader, rawConverter.Converter); err != nil {
		return nil, err
	}

	for i := 0; i < c.Interface.NumMethods(); i++ {
		fun := c.Interface.Method(i)
		def, err := parseMethod(loader, c, fun, rawConverter.Methods[fun.Name()])
		if err != nil {
			return nil, err
		}
		c.Methods[fun.Name()] = def
	}

	return c, nil
}

func parseConverterLines(c *Converter, source string, loader *pkgload.PackageLoader, raw RawLines) error {
	for _, value := range raw.Lines {
		if err := parseConverterLine(c, loader, value); err != nil {
			return formatLineError(raw, source, value, err)
		}
	}

	return nil
}

func parseConverterLine(c *Converter, loader *pkgload.PackageLoader, value string) (err error) {
	cmd, rest := parseCommand(value)
	switch cmd {
	case "converter":
		// only a marker interface
	case "name":
		c.Name, err = parseString(rest)
	case "output:file":
		c.OutputFile, err = parseString(rest)
	case "output:package":
		c.OutputPackageName = ""
		var pkg string
		pkg, err = parseString(rest)

		parts := strings.SplitN(pkg, ":", 2)
		switch len(parts) {
		case 2:
			c.OutputPackageName = parts[1]
			fallthrough
		case 1:
			c.OutputPackagePath = parts[0]
		}
	case "struct:comment":
		c.Comments = append(c.Comments, rest)
	case "extend":
		var methods []*method.Definition
		methods, err = parseExtend(loader, c, strings.Fields(rest))
		c.Extend = append(c.Extend, methods...)
	default:
		_, err = parseCommon(&c.Common, cmd, rest)
	}
	return err
}

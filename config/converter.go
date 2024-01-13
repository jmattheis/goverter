package config

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/pkgload"
)

const (
	configExtend = "extend"
)

type Format string

func (f Format) Struct() bool {
	return f == FormatStruct
}

const (
	FormatStruct    Format = "struct"
	FormatVariables Format = "variables"
)

var DefaultConfigInterface = ConverterConfig{
	OutputFile:        "./generated/generated.go",
	OutputPackageName: "generated",
	OutputFormat:      FormatStruct,
}

var DefaultConfigVariables = ConverterConfig{
	OutputFormat: FormatVariables,
}

type Converter struct {
	ConverterConfig
	Package  string
	FileName string
	typ      types.Type
	Methods  map[string]*Method

	Location string
}

func (c *Converter) requireStruct() error {
	if c.OutputFormat.Struct() {
		return nil
	}
	return fmt.Errorf("not allowed when using goverter:variables")
}

func (c *Converter) IDString() string {
	if c.typ == nil {
		return "var definition"
	}
	return c.typ.String()
}

type ConverterConfig struct {
	Common
	Name              string
	OutputFile        string
	OutputPackagePath string
	OutputPackageName string
	OutputFormat      Format
	Extend            []*method.Definition
	Comments          []string
}

func (conf *ConverterConfig) PackageID() string {
	if conf.OutputPackageName == "" {
		return conf.OutputPackagePath
	}
	return conf.OutputPackagePath + ":" + conf.OutputPackageName
}

func defaultOutputFile(name string) string {
	f := filepath.Base(name)
	ext := filepath.Ext(f)
	return strings.TrimSuffix(f, ext) + ".gen" + ext
}

func parseConverter(loader *pkgload.PackageLoader, rawConverter *RawConverter, global RawLines) (*Converter, error) {
	c, err := initConverter(loader, rawConverter)
	if err != nil {
		return nil, err
	}

	if err := parseConverterLines(c, "global", loader, global); err != nil {
		return nil, err
	}
	if err := parseConverterLines(c, c.IDString(), loader, rawConverter.Converter); err != nil {
		return nil, err
	}

	err = parseMethods(loader, rawConverter, c)
	return c, err
}

func initConverter(loader *pkgload.PackageLoader, rawConverter *RawConverter) (*Converter, error) {
	c := &Converter{
		FileName: rawConverter.FileName,
		Package:  rawConverter.PackagePath,
		Methods:  map[string]*Method{},
		Location: rawConverter.Converter.Location,
	}

	if rawConverter.InterfaceName != "" {
		c.ConverterConfig = DefaultConfigInterface
		v, err := loader.GetOneRaw(c.Package, rawConverter.InterfaceName)
		if err != nil {
			return nil, err
		}

		c.OutputFile = "./generated/generated.go"
		c.OutputPackageName = "generated"
		c.typ = v.Type()
		c.Name = rawConverter.InterfaceName + "Impl"
		c.OutputFormat = FormatStruct
		return c, nil
	}

	c.OutputFormat = FormatVariables
	c.OutputFile = defaultOutputFile(rawConverter.FileName)
	c.OutputPackageName = rawConverter.PackageName
	c.OutputPackagePath = rawConverter.PackagePath
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
	case "converter", "variables":
		// only a marker interface
	case "name":
		if err = c.requireStruct(); err != nil {
			return err
		}
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
		if err = c.requireStruct(); err != nil {
			return err
		}
		c.Comments = append(c.Comments, rest)
	case configExtend:
		for _, name := range strings.Fields(rest) {
			opts := &method.ParseOpts{
				ErrorPrefix:       "error parsing type",
				OutputPackagePath: c.OutputPackagePath,
				Converter:         c.typ,
				Params:            method.ParamsRequired,
			}
			var defs []*method.Definition
			defs, err = loader.GetMatching(c.Package, name, opts)
			if err != nil {
				break
			}
			c.Extend = append(c.Extend, defs...)
		}
	default:
		_, err = parseCommon(&c.Common, cmd, rest)
	}
	return err
}

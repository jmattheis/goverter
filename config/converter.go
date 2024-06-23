package config

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/jmattheis/goverter/enum"
	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/pkgload"
)

const (
	configExtend = "extend"
)

type Format string

const (
	FormatStruct   Format = "struct"
	FormatVariable Format = "assign-variable"
	FormatFunction Format = "function"
)

var DefaultConfigInterface = ConverterConfig{
	OutputFile:        "./generated/generated.go",
	OutputPackageName: "generated",
	Common:            Common{Enum: enum.Config{Enabled: true}},
	OutputFormat:      FormatStruct,
}

var DefaultConfigVariables = ConverterConfig{
	OutputFormat: FormatVariable,
	Common:       Common{Enum: enum.Config{Enabled: true}},
}

type Converter struct {
	ConverterConfig
	Package  string
	FileName string
	typ      types.Type
	Methods  map[string]*Method

	Location string
}

func (c *Converter) typeForMethod() types.Type {
	if c.OutputFormat == FormatFunction {
		return nil
	}
	return c.typ
}

func (c *Converter) requireStruct() error {
	if c.OutputFormat == FormatStruct {
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

func parseConverter(ctx *context, rawConverter *RawConverter, global RawLines) (*Converter, error) {
	c, err := initConverter(ctx.Loader, rawConverter)
	if err != nil {
		return nil, err
	}

	if err := parseConverterLines(ctx, c, "global", global); err != nil {
		return nil, err
	}
	if err := parseConverterLines(ctx, c, c.IDString(), rawConverter.Converter); err != nil {
		return nil, err
	}

	err = parseMethods(ctx, rawConverter, c)
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

	c.OutputFormat = FormatVariable
	c.OutputFile = defaultOutputFile(rawConverter.FileName)
	c.OutputPackageName = rawConverter.PackageName
	c.OutputPackagePath = rawConverter.PackagePath
	return c, nil
}

func parseConverterLines(ctx *context, c *Converter, source string, raw RawLines) error {
	for _, value := range raw.Lines {
		if err := parseConverterLine(ctx, c, value); err != nil {
			return formatLineError(raw, source, value, err)
		}
	}

	return nil
}

func parseConverterLine(ctx *context, c *Converter, value string) (err error) {
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
	case "output:format":
		if len(c.Extend) != 0 {
			return fmt.Errorf("Cannot change output:format after extend functions have been added.\nMove the extend below the output:format setting.")
		}

		c.OutputFormat, err = parseEnum("format", false, rest, FormatFunction, FormatStruct, FormatVariable)
		if err != nil {
			return err
		}

		if c.typ == nil && c.OutputFormat != FormatVariable {
			return fmt.Errorf("unsupported format for goverter:variables")
		}
		if c.typ != nil && c.OutputFormat == FormatVariable {
			return fmt.Errorf("unsupported format for goverter:converter")
		}
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
	case "enum:exclude":
		var pattern enum.IDPattern
		pattern, err = parseIDPattern(c.Package, rest)
		c.Enum.Excludes = append(c.Enum.Excludes, pattern)
	case configExtend:
		for _, name := range strings.Fields(rest) {
			opts := &method.ParseOpts{
				ErrorPrefix:       "error parsing type",
				OutputPackagePath: c.OutputPackagePath,
				Converter:         c.typeForMethod(),
				Params:            method.ParamsRequired,
			}
			var defs []*method.Definition
			defs, err = ctx.Loader.GetMatching(c.Package, name, opts)
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

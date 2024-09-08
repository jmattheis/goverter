package config

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"

	"github.com/jmattheis/goverter/method"
)

const (
	configMap     = "map"
	configDefault = "default"
)

var StructMethodContextRegex = regexp.MustCompile(".*")

type Method struct {
	*method.Definition
	Common

	Constructor *method.Definition
	AutoMap     []string
	Fields      map[string]*FieldMapping
	EnumMapping *EnumMapping

	RawFieldSettings []string

	Location string
}

type FieldMapping struct {
	Source   string
	Function *method.Definition
	Ignore   bool
}

func (m *Method) Field(targetName string) *FieldMapping {
	target, ok := m.Fields[targetName]
	if !ok {
		target = &FieldMapping{}
		m.Fields[targetName] = target
	}
	return target
}

func parseMethods(ctx *context, rawConverter *RawConverter, c *Converter) error {
	if c.typ != nil {
		interf := c.typ.Underlying().(*types.Interface)
		for i := 0; i < interf.NumMethods(); i++ {
			fun := interf.Method(i)
			def, err := parseMethod(ctx, c, fun, rawConverter.Methods[fun.Name()])
			if err != nil {
				return err
			}
			c.Methods = append(c.Methods, def)
		}
		return nil
	}
	for name, lines := range rawConverter.Methods {
		fun, err := ctx.Loader.GetOneRaw(c.Package, name)
		if err != nil {
			return err
		}
		def, err := parseMethod(ctx, c, fun, lines)
		if err != nil {
			return err
		}
		c.Methods = append(c.Methods, def)
	}
	return nil
}

func parseMethod(ctx *context, c *Converter, obj types.Object, rawMethod RawLines) (*Method, error) {
	m := &Method{
		Common:      c.Common,
		Fields:      map[string]*FieldMapping{},
		Location:    rawMethod.Location,
		EnumMapping: &EnumMapping{Map: map[string]string{}},
	}

	for _, value := range rawMethod.Lines {
		if err := parseMethodLine(ctx, c, m, value); err != nil {
			return m, formatLineError(rawMethod, obj.String(), value, err)
		}
	}

	def, err := method.Parse(obj, &method.ParseOpts{
		ErrorPrefix:       "error parsing converter method",
		Location:          rawMethod.Location,
		Converter:         nil,
		OutputPackagePath: c.OutputPackagePath,
		Params:            method.ParamsRequired,
		ContextMatch:      m.ArgContextRegex,
		Generated:         true,
	})

	m.Definition = def

	return m, err
}

func parseMethodLine(ctx *context, c *Converter, m *Method, value string) (err error) {
	cmd, rest := parseCommand(value)
	fieldSetting := false
	switch cmd {
	case configMap:
		fieldSetting = true
		var source, target, custom string
		source, target, custom, err = parseMethodMap(rest)
		if err != nil {
			return err
		}
		f := m.Field(target)
		f.Source = source

		if custom != "" {
			opts := &method.ParseOpts{
				ErrorPrefix:       "error parsing type",
				OutputPackagePath: c.OutputPackagePath,
				Converter:         c.typeForMethod(),
				Params:            method.ParamsOptional,
				AllowTypeParams:   true,
				ContextMatch:      m.ArgContextRegex,
			}
			f.Function, err = ctx.Loader.GetOne(c.Package, custom, opts)
		}
	case "ignore":
		fieldSetting = true
		fields := strings.Fields(rest)
		for _, f := range fields {
			m.Field(f).Ignore = true
		}
	case "enum:map":
		fields := strings.Fields(rest)
		if len(fields) != 2 {
			return fmt.Errorf("invalid fields")
		}

		if IsEnumAction(fields[1]) {
			err = validateEnumAction(fields[1])
		}

		m.EnumMapping.Map[fields[0]] = fields[1]
	case "enum:transform":
		fields := strings.SplitN(rest, " ", 2)

		config := ""
		if len(fields) == 2 {
			config = fields[1]
		}

		var t ConfiguredTransformer
		t, err = parseTransformer(ctx, fields[0], config)
		m.EnumMapping.Transformers = append(m.EnumMapping.Transformers, t)
	case "autoMap":
		fieldSetting = true
		var s string
		s, err = parseString(rest)
		m.AutoMap = append(m.AutoMap, strings.TrimSpace(s))
	case configDefault:
		opts := &method.ParseOpts{
			ErrorPrefix:       "error parsing type",
			OutputPackagePath: c.OutputPackagePath,
			Converter:         c.typeForMethod(),
			Params:            method.ParamsOptional,
			AllowTypeParams:   true,
			ContextMatch:      m.ArgContextRegex,
		}
		m.Constructor, err = ctx.Loader.GetOne(c.Package, rest, opts)
	default:
		fieldSetting, err = parseCommon(&m.Common, cmd, rest)
	}
	if fieldSetting {
		m.RawFieldSettings = append(m.RawFieldSettings, value)
	}
	return err
}

func parseMethodMap(remaining string) (source, target, custom string, err error) {
	parts := strings.SplitN(remaining, "|", 2)
	if len(parts) == 2 {
		custom = strings.TrimSpace(parts[1])
	}

	fields := strings.Fields(parts[0])
	switch len(fields) {
	case 1:
		target = fields[0]
	case 2:
		source = fields[0]
		target = fields[1]
	case 0:
		err = fmt.Errorf("missing target field")
	default:
		err = fmt.Errorf("too many fields expected at most 2 fields got %d: %s", len(fields), remaining)
	}
	if err == nil && strings.ContainsRune(target, '.') {
		err = fmt.Errorf("the mapping target %q must be a field name but was a path.\nDots \".\" are not allowed.", target)
	}
	return source, target, custom, err
}

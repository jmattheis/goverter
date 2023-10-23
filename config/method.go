package config

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/pkgload"
)

type Method struct {
	*method.Definition
	Common
	AutoMap []string
	Fields  map[string]*FieldMapping
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

func parseMethod(loader *pkgload.PackageLoader, c *Converter, fn *types.Func, rawMethod RawLines) (*Method, error) {
	def, err := method.Parse(&method.ParseOpts{
		Obj:          fn,
		ErrorPrefix:  "error parsing converter method",
		Converter:    nil,
		EmptySource:  false,
		ConvFunction: true,
	})
	if err != nil {
		return nil, err
	}

	m := &Method{
		Definition: def,
		Common:     c.Common,
		Fields:     map[string]*FieldMapping{},
	}

	for _, value := range rawMethod.Lines {
		if err := parseMethodLine(loader, c, m, value); err != nil {
			return m, formatLineError(rawMethod, fn.String(), value, err) // TODO get method type
		}
	}
	return m, nil
}

func parseMethodLine(loader *pkgload.PackageLoader, c *Converter, m *Method, value string) (err error) {
	cmd, rest := parseCommand(value)
	switch cmd {
	case "map":
		err = parseMethodMap(loader, c, m, rest)
	case "ignore":
		fields := strings.Fields(rest)
		for _, f := range fields {
			m.Field(f).Ignore = true
		}
	case "autoMap":
		var s string
		s, err = parseString(rest)
		m.AutoMap = append(m.AutoMap, strings.TrimSpace(s))
	default:
		err = parseCommon(&m.Common, cmd, rest)
	}
	return err
}

func parseMethodMap(loader *pkgload.PackageLoader, c *Converter, m *Method, remaining string) (err error) {
	parts := strings.SplitN(remaining, "|", 2)
	fields := strings.Fields(parts[0])
	custom := ""
	if len(parts) == 2 {
		custom = strings.TrimSpace(parts[1])
	}

	switch len(fields) {
	case 1:
		if custom != "" {
			m.Field(fields[0]).Function, err = parseMapExtend(loader, c, custom)
		}
	case 2:
		f := m.Field(fields[1])
		f.Source = fields[0]
		if custom != "" {
			f.Function, err = parseMapExtend(loader, c, custom)
		}
	case 0:
		err = fmt.Errorf("missing target field")
	default:
		err = fmt.Errorf("too many fields expected at most 2 fields got %d: %s", len(fields), remaining)
	}
	return err
}

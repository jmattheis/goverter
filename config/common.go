package config

import (
	"fmt"
	"regexp"

	"github.com/jmattheis/goverter/config/parse"
	"github.com/jmattheis/goverter/enum"
)

type Common struct {
	FieldSettings                      []string
	WrapErrors                         bool
	WrapErrorsUsing                    string
	IgnoreUnexported                   bool
	IgnoreBasicZeroValueField          bool
	IgnoreStructZeroValueField         bool
	IgnoreNillableZeroValueField       bool
	MatchIgnoreCase                    bool
	IgnoreMissing                      bool
	SkipCopySameType                   bool
	UseZeroValueOnPointerInconsistency bool
	UseUnderlyingTypeMethods           bool
	DefaultUpdate                      bool
	ArgContextRegex                    *regexp.Regexp
	Enum                               enum.Config
}

func parseCommon(c *Common, cmd, rest string) (fieldSetting bool, err error) {
	switch cmd {
	case "wrapErrors":
		if c.WrapErrorsUsing != "" {
			return false, fmt.Errorf("cannot be used in combination with wrapErrorsUsing")
		}
		c.WrapErrors, err = parse.Bool(rest)
	case "wrapErrorsUsing":
		if c.WrapErrors {
			return false, fmt.Errorf("cannot be used in combination with wrapErrors")
		}
		c.WrapErrorsUsing, err = parse.String(rest)
	case "ignoreUnexported":
		fieldSetting = true
		c.IgnoreUnexported, err = parse.Bool(rest)
	case "update:ignoreZeroValueField":
		fieldSetting = true
		c.IgnoreBasicZeroValueField, err = parse.Bool(rest)
		c.IgnoreStructZeroValueField = c.IgnoreBasicZeroValueField
		c.IgnoreNillableZeroValueField = c.IgnoreBasicZeroValueField
	case "update:ignoreZeroValueField:basic":
		c.IgnoreBasicZeroValueField, err = parse.Bool(rest)
	case "update:ignoreZeroValueField:struct":
		c.IgnoreStructZeroValueField, err = parse.Bool(rest)
	case "update:ignoreZeroValueField:nillable":
		c.IgnoreNillableZeroValueField, err = parse.Bool(rest)
	case "default:update":
		c.DefaultUpdate, err = parse.Bool(rest)
	case "matchIgnoreCase":
		fieldSetting = true
		c.MatchIgnoreCase, err = parse.Bool(rest)
	case "ignoreMissing":
		fieldSetting = true
		c.IgnoreMissing, err = parse.Bool(rest)
	case "skipCopySameType":
		c.SkipCopySameType, err = parse.Bool(rest)
	case "useZeroValueOnPointerInconsistency":
		c.UseZeroValueOnPointerInconsistency, err = parse.Bool(rest)
	case "useUnderlyingTypeMethods":
		c.UseUnderlyingTypeMethods, err = parse.Bool(rest)
	case "enum":
		c.Enum.Enabled, err = parse.Bool(rest)
	case "arg:context:regex":
		c.ArgContextRegex, err = parse.Regex(rest)
	case "enum:unknown":
		c.Enum.Unknown, err = parse.String(rest)
		if err == nil && IsEnumAction(c.Enum.Unknown) {
			err = validateEnumAction(c.Enum.Unknown)
		}
	case "":
		err = fmt.Errorf("missing setting key")
	default:
		err = fmt.Errorf("unknown setting: %s", cmd)
	}

	return fieldSetting, err
}

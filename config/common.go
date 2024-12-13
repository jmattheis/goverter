package config

import (
	"fmt"
	"regexp"

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
		c.WrapErrors, err = parseBool(rest)
	case "wrapErrorsUsing":
		if c.WrapErrors {
			return false, fmt.Errorf("cannot be used in combination with wrapErrors")
		}
		c.WrapErrorsUsing, err = parseString(rest)
	case "ignoreUnexported":
		fieldSetting = true
		c.IgnoreUnexported, err = parseBool(rest)
	case "update:ignoreZeroValueField":
		fieldSetting = true
		c.IgnoreBasicZeroValueField, err = parseBool(rest)
		c.IgnoreStructZeroValueField = c.IgnoreBasicZeroValueField
		c.IgnoreNillableZeroValueField = c.IgnoreBasicZeroValueField
	case "update:ignoreZeroValueField:basic":
		c.IgnoreBasicZeroValueField, err = parseBool(rest)
	case "update:ignoreZeroValueField:struct":
		c.IgnoreStructZeroValueField, err = parseBool(rest)
	case "update:ignoreZeroValueField:nillable":
		c.IgnoreNillableZeroValueField, err = parseBool(rest)
	case "default:update":
		c.DefaultUpdate, err = parseBool(rest)
	case "matchIgnoreCase":
		fieldSetting = true
		c.MatchIgnoreCase, err = parseBool(rest)
	case "ignoreMissing":
		fieldSetting = true
		c.IgnoreMissing, err = parseBool(rest)
	case "skipCopySameType":
		c.SkipCopySameType, err = parseBool(rest)
	case "useZeroValueOnPointerInconsistency":
		c.UseZeroValueOnPointerInconsistency, err = parseBool(rest)
	case "useUnderlyingTypeMethods":
		c.UseUnderlyingTypeMethods, err = parseBool(rest)
	case "enum":
		c.Enum.Enabled, err = parseBool(rest)
	case "arg:context:regex":
		c.ArgContextRegex, err = parseRegex(rest)
	case "enum:unknown":
		c.Enum.Unknown, err = parseString(rest)
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

package config

import "fmt"

type Common struct {
	FieldSettings                      []string
	WrapErrors                         bool
	WrapErrorsUsing                    string
	IgnoreUnexported                   bool
	MatchIgnoreCase                    bool
	IgnoreMissing                      bool
	SkipCopySameType                   bool
	UseZeroValueOnPointerInconsistency bool
	UseUnderlyingTypeMethods           bool
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
	case "":
		err = fmt.Errorf("missing setting key")
	default:
		err = fmt.Errorf("unknown setting: %s", cmd)
	}

	return fieldSetting, err
}

package config

import "fmt"

type Common struct {
	WrapErrors                         bool
	IgnoreUnexported                   bool
	MatchIgnoreCase                    bool
	IgnoreMissing                      bool
	SkipCopySameType                   bool
	UseZeroValueOnPointerInconsistency bool
}

func parseCommon(c *Common, cmd, rest string) (err error) {
	switch cmd {
	case "wrapErrors":
		c.WrapErrors, err = parseBool(rest)
	case "ignoreUnexported":
		c.IgnoreUnexported, err = parseBool(rest)
	case "matchIgnoreCase":
		c.MatchIgnoreCase, err = parseBool(rest)
	case "ignoreMissing":
		c.IgnoreMissing, err = parseBool(rest)
	case "skipCopySameType":
		c.SkipCopySameType, err = parseBool(rest)
	case "useZeroValueOnPointerInconsistency":
		c.UseZeroValueOnPointerInconsistency, err = parseBool(rest)
	case "":
		err = fmt.Errorf("missing setting key")
	default:
		err = fmt.Errorf("unknown setting: %s", cmd)
	}

	return
}

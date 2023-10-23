package config

import (
	"fmt"
	"strings"
)

func parseCommand(value string) (string, string) {
	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func parseBool(remaining string) (bool, error) {
	fields := strings.Fields(remaining)
	switch {
	case len(fields) == 0:
		return true, nil
	case len(fields) == 1:
		switch fields[0] {
		case "yes":
			return true, nil
		case "no":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean value: '%s' must be one of 'yes', 'no'", fields[0])
		}
	default:
		return false, fmt.Errorf("expected at most one value but got %d: %#v", len(fields), fields)
	}
}

func parseString(remaining string) (string, error) {
	fields := strings.Fields(remaining)
	if len(fields) != 1 {
		return "", fmt.Errorf("must have one value but got %d: %#v", len(fields), remaining)
	}
	return fields[0], nil
}

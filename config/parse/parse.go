package parse

import (
	"fmt"
	"regexp"
	"strings"
)

func Command(value string) (string, string) {
	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func Enum[T ~string](empty bool, remaining string, values ...T) (T, error) {
	fields := strings.Fields(remaining)

	switch {
	case len(fields) == 0 && empty:
		return "", nil
	case len(fields) == 1:
		for _, value := range values {
			if fields[0] == string(value) {
				return value, nil
			}
		}

		return "", fmt.Errorf("invalid value: '%s' must be one of: %s", fields[0], formatValues(values))
	default:
		return "", fmt.Errorf("invalid value: expected one value but got %d: %s", len(fields), fields)
	}
}

func formatValues[T ~string](values []T) string {
	strs := make([]string, len(values))
	for i, id := range values {
		strs[i] = string(id)
	}
	return strings.Join(strs, ", ")
}

func Bool(remaining string) (bool, error) {
	val, err := Enum(true, remaining, "yes", "no")
	return val == "" || val == "yes", err
}

func String(remaining string) (string, error) {
	fields := strings.Fields(remaining)
	if len(fields) != 1 {
		return "", fmt.Errorf("must have one value but got %d: %#v", len(fields), remaining)
	}
	return fields[0], nil
}

func Regex(remaining string) (*regexp.Regexp, error) {
	value, err := String(remaining)
	if err != nil {
		return nil, err
	}
	return regexp.Compile(value)
}

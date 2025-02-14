package parse

import (
	"path/filepath"
	"strings"
)

func File(cwd, rest string) (string, error) {
	field, err := String(rest)
	if err != nil {
		return field, err
	}

	if strings.HasPrefix(field, "@cwd/") {
		return filepath.Abs(filepath.Join(cwd, strings.TrimPrefix(field, "@cwd/")))
	}
	return field, nil
}

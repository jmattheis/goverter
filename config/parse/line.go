package parse

import (
	"bufio"
	"strings"
)

const (
	Prefix    = "goverter"
	Delimiter = ":"
)

func SettingLines(comment string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(comment))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, Prefix+Delimiter) {
			line := strings.TrimPrefix(line, Prefix+Delimiter)
			lines = append(lines, line)
		}
	}
	return lines
}

package example

import "time"

// goverter:converter
// goverter:extend FormatTime
type Converter interface {
	// goverter:context dateFormat
	Convert(source map[string]Input, dateFormat string) map[string]Output
}

// goverter:context dateFormat
func FormatTime(t time.Time, dateFormat string) string {
	return t.Format(dateFormat)
}

type Input struct {
	Name      string
	CreatedAt time.Time
}
type Output struct {
	Name      string
	CreatedAt string
}

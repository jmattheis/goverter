package example

import "time"

// goverter:converter
// goverter:arg:context:regex .+Format
// goverter:extend FormatTime
type Converter interface {
	Convert(source map[string]Input, dateFormat string) map[string]Output
}

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

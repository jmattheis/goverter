package example

import "time"

// goverter:converter
// goverter:extend FormatTime
type Converter interface {
	Convert(source map[string]Input, ctxFormat string) map[string]Output
}

func FormatTime(t time.Time, ctxFormat string) string {
	return t.Format(ctxFormat)
}

type Input struct {
	Name      string
	CreatedAt time.Time
}
type Output struct {
	Name      string
	CreatedAt string
}

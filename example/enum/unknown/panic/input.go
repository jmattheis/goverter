package example

import (
	"github.com/jmattheis/goverter/example/enum/unknown/input"
	"github.com/jmattheis/goverter/example/enum/unknown/output"
)

// goverter:converter
// goverter:enum:unknown @panic
type Converter interface {
	Convert(input.Color) output.Color
}

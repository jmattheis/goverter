package example

import (
	"github.com/jmattheis/goverter/example/enum/unknown/input"
	"github.com/jmattheis/goverter/example/enum/unknown/output"
)

// goverter:converter
// goverter:enum:unknown @error
type Converter interface {
	Convert(input.Color) (output.Color, error)
}

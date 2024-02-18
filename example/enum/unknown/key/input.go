package example

import (
	"github.com/jmattheis/goverter/example/enum/unknown/key/input"
	"github.com/jmattheis/goverter/example/enum/unknown/key/output"
)

// goverter:converter
type Converter interface {
    // goverter:enum:unknown Unknown
    Convert(input.Color) output.Color
}

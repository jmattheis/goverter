package example

import (
	"github.com/jmattheis/goverter/example/enum/map/input"
	"github.com/jmattheis/goverter/example/enum/map/output"
)

// goverter:converter
// goverter:enum:unknown @panic
type Converter interface {
    // goverter:enum:map Gray Grey
    Convert(input.Color) output.Color
}

package example

// goverter:converter
// goverter:enum:unknown @panic
type Converter interface {
	// goverter:enum:transform trim-prefix Color
	Convert(InputColor) OutputColor
}

type InputColor int

const (
	ColorGreen InputColor = iota
	ColorBlue
	ColorRed
)

type OutputColor string

const (
	Green OutputColor = "green"
	Blue  OutputColor = "blue"
	Red   OutputColor = "red"
)

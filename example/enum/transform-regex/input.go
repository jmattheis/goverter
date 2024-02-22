package example

// goverter:converter
// goverter:enum:unknown @panic
type Converter interface {
	// goverter:enum:transform regex Col(\w+) ${1}Color
	Convert(InputColor) OutputColor
	// goverter:enum:transform regex (\w+)Color Col${1}
	Convert2(OutputColor) InputColor
}

type InputColor int

const (
	ColGreen InputColor = iota
	ColBlue
	ColRed
)

type OutputColor string

const (
	GreenColor OutputColor = "green"
	BlueColor  OutputColor = "blue"
	RedColor   OutputColor = "red"
)

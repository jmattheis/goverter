package update

// goverter:converter
type Converter interface {
	// goverter:update target
	Convert(source Input, target *Output)
}

type Input struct {
	Name    *string
	Aliases []string
	Age     int
}

type Output struct {
	Name    *string
	Aliases []string
	Age     int
}

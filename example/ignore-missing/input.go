package example

// goverter:converter
type Converter interface {
	// goverter:ignoreMissing
	Convert(source Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	Name   string
	Age    int
	Street string
}

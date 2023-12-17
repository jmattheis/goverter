package example

// goverter:converter
type Converter interface {
	// goverter:ignore Age
	Convert(source Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	Name string
	Age  int
}

package example

// goverter:converter
type Converter interface {
	// goverter:ignoreUnexported
	Convert(source Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	Name string
	// goverter will skip this field
	age    int
	street string
}

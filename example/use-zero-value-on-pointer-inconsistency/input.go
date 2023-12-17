package example

// goverter:converter
type Converter interface {
	// goverter:useZeroValueOnPointerInconsistency
	Convert(source Input) Output
}

type Input struct {
	Name *string
	Age  int
}

type Output struct {
	Name string
	Age  int
}

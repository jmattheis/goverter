package example

// goverter:converter
type Converter interface {
	// goverter:map Nested.LastName Surname
	Convert(Input) Output
}

type Input struct {
	Age    int
	Nested NestedInput
}
type NestedInput struct {
	LastName string
}
type Output struct {
	Age     int
	Surname string
}

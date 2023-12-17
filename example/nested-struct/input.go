package example

// goverter:converter
type Converter interface {
	Convert([]Input) []Output

	// goverter:map LastName Surname
	ConvertNested(NestedInput) NestedOutput
}

type Input struct {
	Name   string
	Nested NestedInput
}
type NestedInput struct {
	LastName string
}

type Output struct {
	Name   string
	Nested NestedOutput
}
type NestedOutput struct {
	Surname string
}

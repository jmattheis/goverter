package example

// goverter:converter
type Converter interface {
	Convert(Input) Output

	// goverter:map LastName Surname
	ConvertNested(NestedInput) NestedOutput
}

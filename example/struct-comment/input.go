package example

// goverter:converter
// goverter:struct:comment // MyConverterImpl
// goverter:struct:comment //
// goverter:struct:comment // More detailed
type MultipleSingleLine interface {
	Convert(Input) Output
}

// goverter:converter
// goverter:struct:comment single comment
type SingleComment interface {
	Convert(Input) Output
}

// goverter:converter
// goverter:struct:comment MyConverterImpl
// goverter:struct:comment
// goverter:struct:comment More detailed
type MultiLine interface {
	Convert(Input) Output
}

type Input struct{ Name string }
type Output struct{ Name string }

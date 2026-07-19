package example

// goverter:converter
// goverter:annotate:unmapped
type Converter interface {
	// goverter:ignore Age
	// goverter:ignoreMissing
	// goverter:ignoreUnexported
	Convert(source Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	Name       string
	Age        int
	Missing    string
	unexported string
}

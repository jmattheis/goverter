package example

// goverter:converter
type Converter interface {
	// goverter:matchIgnoreCase
	Convert(Input) Output
}

type Input struct {
	Age      int
	Fullname string
}
type Output struct {
	Age      int
	FULLNAME string
}

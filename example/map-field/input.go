package example

// goverter:converter
type Converter interface {
	// goverter:map LastName Surname
	Convert(Input) Output
}

type Input struct {
	Age      int
	LastName string
}
type Output struct {
	Age     int
	Surname string
}

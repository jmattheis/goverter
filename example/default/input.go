package example

// goverter:converter
type Converter interface {
	// goverter:default NewOutput
	Convert(*Input) *Output
}

type Input struct {
	Age  int
	Name *string
}
type Output struct {
	Age  int
	Name *string
}

func NewOutput() *Output {
	name := "jmattheis"
	return &Output{
		Age:  42,
		Name: &name,
	}
}

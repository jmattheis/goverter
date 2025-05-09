package example

// goverter:converter
// goverter:setters:enabled
type Converter interface {
	// goverter:ignore name
	Convert(Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	name string
}

func (o *Output) SetName(name string) {
	o.name = name
}

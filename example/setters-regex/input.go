package example

// goverter:converter
// goverter:setters:enabled
// goverter:setters:regex With([A-Z].*)
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

func (o *Output) WithName(name string) {
	o.name = name
}

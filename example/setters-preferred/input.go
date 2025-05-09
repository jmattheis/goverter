package example

// goverter:converter
// goverter:setters:enabled
// goverter:setters:preferred
type Converter interface {
	Convert(Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	Name string
}

func (o *Output) SetName(name string) {
	o.Name = name
}

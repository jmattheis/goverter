package example

// goverter:converter
// goverter:getters:enabled
type Converter interface {
	Convert(Input) Output
}

type Input struct {
	name string
}

func (i Input) GetName() string {
	return i.name
}

type Output struct {
	Name string
}

package example

// goverter:converter
// goverter:getters:enabled
// goverter:getters:preferred
type Converter interface {
	Convert(Input) Output
}

type Input struct {
	Name string
}

func (i Input) GetName() string {
	return i.Name
}

type Output struct {
	Name string
}

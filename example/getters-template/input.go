package example

// goverter:converter
// goverter:getters:enabled
// goverter:getters:template This{{.}}
type Converter interface {
	Convert(Input) Output
}

type Input struct {
	Name string
}

func (i Input) ThisName() string {
	return i.Name
}

type Output struct {
	Name string
}

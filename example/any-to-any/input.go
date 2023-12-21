package example

// goverter:converter
// goverter:extend ConvertAny
type Converter interface {
	Convert(Input) Output
}

func ConvertAny(value interface{}) interface{} {
	return value
}

type Input struct{ Value interface{} }
type Output struct{ Value interface{} }

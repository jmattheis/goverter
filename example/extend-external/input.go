package example

// goverter:converter
// goverter:extend strconv:Atoi
type Converter interface {
	Convert(Input) (Output, error)
}

type Input struct{ Value string }
type Output struct{ Value int }

package example

// goverter:converter
// goverter:extend strconv:Atoi
// goverter:wrapErrors
type Converter interface {
	Convert(source Input) (Output, error)
}

type Input struct {
	PostalCode string
}
type Output struct {
	PostalCode int
}

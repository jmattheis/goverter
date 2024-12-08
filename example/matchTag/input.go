package example

// goverter:converter
// goverter:matchTag json
type Converter interface {
	Convert(Input) Output
}

type Input struct {
	Name string `json:"blame"`
}
type Output struct {
	Game string `json:"blame"`
}

package example

// goverter:converter
// goverter:name RenamedConverter
type Converter interface {
	Convert(Input) Output
}

type Input struct {
	Name string
}
type Output struct {
	Name string
}

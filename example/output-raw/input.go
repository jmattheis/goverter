package raw

// goverter:converter
// goverter:output:raw func Hello() string {
// goverter:output:raw    return "World!"
// goverter:output:raw }
type Converter interface {
	Convert(source Input) Output
}

type Input struct{ Name string }
type Output struct{ Name string }

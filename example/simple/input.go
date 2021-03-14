//go:generate go run github.com/jmattheis/goverter/cmd/goverter github.com/jmattheis/goverter/example/simple
package simple

// goverter:converter
type Converter interface {
	Convert(source []Input) []Output
}

type Input struct {
	Name string
	Age  int
}
type Output struct {
	Name string
	Age  int
}

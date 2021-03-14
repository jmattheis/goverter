//go:generate go run github.com/jmattheis/go-genconv/cmd/go-genconv github.com/jmattheis/go-genconv/example/simple
package simple

// genconv:converter
type Converter interface {
	Convert(source []Input) []Output
}

type Input struct {
	Name string
	Age int
}
type Output struct {
	Name string
	Age int
}

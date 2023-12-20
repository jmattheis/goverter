package samepackage

// goverter:converter
// goverter:output:file ./generated.go
// goverter:output:package github.com/jmattheis/goverter/example/samepackage
type Converter interface {
	Convert(source *Input) *Output
}

type Input struct { Name string }
type Output struct { Name string }

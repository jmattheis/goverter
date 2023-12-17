package example

// goverter:converter
// goverter:output:file ./a/generated.go
// goverter:output:package goverter/example/a
type RootA interface {
	Convert([]bool) []bool
}

// goverter:converter
// goverter:output:file ./b/generated.go
// goverter:output:package goverter/example/b
type RootB interface {
	Convert([]string) []string
}

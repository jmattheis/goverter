package c

// goverter:converter
// goverter:output:file ../a/generated.go
// goverter:output:package goverter/example/a
type CIntoA interface {
	Convert([]int) []int
}

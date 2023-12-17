package example

import "fmt"

// goverter:converter
// goverter:extend IntToString
type Converter interface {
	Convert(Input) Output
}
type Input struct {
	Name string
	Age  int
}
type Output struct {
	Name string
	Age  string
}

func IntToString(i int) string {
	return fmt.Sprint(i)
}

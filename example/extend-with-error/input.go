package example

import "strconv"

// goverter:converter
// goverter:extend StringToInt
type Converter interface {
	Convert(Input) (Output, error)
}

type Input struct{ Value string }
type Output struct{ Value int }

func StringToInt(value string) (int, error) {
	i, err := strconv.Atoi(value)
	return i, err
}

package simple

import "time"

// goverter:converter
// goverter:extend ConvertTime
type Converter interface {
	Convert(source Input) Output
}

// It is possible to only call the ConvertTime function
// for one specific conversion method like this

// goverter:converter
type ConverterLocally interface {
	// goverter:map CreatedAt | ConvertTime
	Convert(source Input) Output
}

func ConvertTime(t time.Time) time.Time {
	return t
}

type Input struct {
	Name      string
	CreatedAt time.Time
}

type Output struct {
	Name      string
	CreatedAt time.Time
}

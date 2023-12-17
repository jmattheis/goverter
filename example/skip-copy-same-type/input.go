package example

// goverter:converter
// goverter:skipCopySameType
type Converter interface {
	Convert(source Input) Output
}

type Input struct {
	Name       *string
	ItemCounts map[string]int
}

type Output struct {
	Name       *string
	ItemCounts map[string]int
}

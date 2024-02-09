package example

// goverter:converter
// goverter:wrapErrorsUsing github.com/goverter/patherr
// goverter:extend strconv:Atoi
type Converter interface {
	Convert(source map[int]Input) (map[int]Output, error)
}

type Input struct {
	Friends    []Input
	Age        string
	Attributes map[string]string
}

type Output struct {
	Friends    []Output
	Age        int
	Attributes map[string]int
}

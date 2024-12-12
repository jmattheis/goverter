package example

// goverter:converter
// goverter:output:package "github.com/jmattheis/goverter/example/default/generated"
type Converter interface {
	// goverter:default NewOutput
	// goverter:ignore Age
	ConvertInterfaceStruct(input Input) Output
}

// goverter:variables
// goverter:output:file generated/generated.go
// goverter:output:package "github.com/jmattheis/goverter/example/default/generated"
var (
	// goverter:default NewOutputPointer
	// goverter:ignore Age
	ConvertVarPointer func(input Input) *Output
)

type Input struct {
	Age  int
	Name *string
}
type Output struct {
	Age  int
	Name *string
}

func NewOutput(input Input) Output {
	return Output{
		Age:  42,
		Name: input.Name,
	}
}

func NewOutputPointer(input *Input) *Output {
	return &Output{
		Age:  42,
		Name: input.Name,
	}
}

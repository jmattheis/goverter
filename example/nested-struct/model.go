package example

type Input struct {
	Name   string
	Nested NestedInput
}
type Output struct {
	Name   string
	Nested NestedOutput
}

type NestedInput struct{ LastName string }
type NestedOutput struct{ Surname string }

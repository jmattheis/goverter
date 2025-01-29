package package_name_that_doesnt_match_dir_name

// goverter:converter
// goverter:output:file ./generated.go
type Converter interface {
	Convert(source *Input) *Output
}

type Input struct{ Name string }
type Output struct{ Name string }
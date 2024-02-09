package example

// goverter:converter
// goverter:wrapErrorsUsing goverter/example/patherr
// goverter:output:file ./generated/minimal.go
// goverter:extend strconv:Atoi
type Minimal interface {
	Convert(source map[int]Input) (map[int]Output, error)
}

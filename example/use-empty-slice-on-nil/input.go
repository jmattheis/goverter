package emptyonnil

// goverter:converter
// goverter:useEmptySliceOnNil
// goverter:useEmptyMapOnNil
type Converter interface {
	Convert(source []Input) []Output
	ConvertMap(source map[string]string) map[string]string
}

// goverter:converter
type StrictConverter interface {
	Convert(source []Input) []Output
	ConvertMap(source map[string]string) map[string]string
}

type Input struct {
	Name string
}

type Output struct {
	Name string
}

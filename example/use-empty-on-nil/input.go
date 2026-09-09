package emptyonnil

// goverter:converter
// goverter:useEmptyOnNil
type Converter interface {
	ConvertSlice(source []string) []string
	ConvertMap(source map[string]string) map[string]string
}

// goverter:converter
type StrictConverter interface {
	ConvertSlice(source []string) []string
	ConvertMap(source map[string]string) map[string]string
}

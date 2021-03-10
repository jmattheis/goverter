//go:generate go run github.com/jmattheis/go-genconv/cmd/go-genconv github.com/jmattheis/go-genconv/test/primitive
package primitive

// genconv:converter
type Converter interface {
	ConvertString(source string) string
	ConvertInt(source int) int
	ConvertInt8(source int8) int8
	ConvertInt16(source int16) int16
}

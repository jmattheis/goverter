//go:generate go run github.com/jmattheis/go-genconv/cmd/go-genconv github.com/jmattheis/go-genconv/test/comments
package comments

// genconv:converter
type Converter interface {
	ConvertString(source string) string
}

// genconv:converter
// genconv:name RenamedConverter
type (
	Converter2 interface {
		ConvertString(source string) string
	}
)

type (
	// genconv:converter
	Converter3 interface {
		ConvertString(source string) string
	}
	// genconv:converter
	Converter4 interface {
		ConvertString(source string) string
	}
)

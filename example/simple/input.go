//go:generate go run github.com/jmattheis/go-genconv/cmd/go-genconv github.com/jmattheis/go-genconv/example/simple
package simple

// genconv:converter
type (
	Converter interface {
		ConvertModel(DBModel) ExternalModel
	}
)

type DBModel struct {
	Name string
	Age  int
}

type ExternalModel struct {
	Name string
	Age  int
}

//go:generate go run github.com/jmattheis/go-genconv/cmd/go-genconv github.com/jmattheis/go-genconv/example/errors
package errors

import (
	"fmt"
	"math/big"
)

// genconv:converter
type Converter interface {
	ToExternal(DBModel) ExternalModel
	ToDB(model ExternalModel) (DBModel, error)

	// genconv:delegate BigIntToString
	BigIntToString(big.Int) string
	// genconv:delegate StringToBigInt
	StringToBigInt(string) (big.Int, error)
}

func StringToBigInt(_ Converter, value string) (big.Int, error) {
	i := big.Int{}
	_, ok := i.SetString(value, 10)
	if !ok {
		return i, fmt.Errorf("could not convert string to big.Int")
	}
	return i, nil
}
func BigIntToString(_ Converter, value big.Int) string {
	return value.String()
}

type DBModel struct {
	ID big.Int

	Name string
	Age  int
}

type ExternalModel struct {
	ID   string
	Name string
	Age  int
}

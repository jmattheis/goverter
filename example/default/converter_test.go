package example_test

import (
	"testing"

	default1 "github.com/jmattheis/goverter/example/default"
	"github.com/jmattheis/goverter/example/default/generated"

	"github.com/stretchr/testify/require"
)

func TestConvertInterfaceSuccess(t *testing.T) {
	var c default1.Converter = &generated.ConverterImpl{}

	input := default1.Input{
		Age:  20,
		Name: p("tester"),
	}

	output := c.ConvertInterfaceStruct(input)

	expected := default1.Output{
		Age:  42,
		Name: p("tester"),
	}

	require.Equal(t, expected, output)
}

func TestConvertVarSuccess(t *testing.T) {
	input := default1.Input{
		Age:  20,
		Name: p("tester"),
	}

	output := default1.ConvertVarPointer(input)

	expected := default1.Output{
		Age:  42,
		Name: p("tester"),
	}

	require.Equal(t, expected, *output)
}

func p(s string) *string {
	return &s
}

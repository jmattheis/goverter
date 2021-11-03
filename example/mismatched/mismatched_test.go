package mismatched_test

import (
	"testing"

	"github.com/jmattheis/goverter/example/mismatched"
	"github.com/jmattheis/goverter/example/mismatched/generated"
	"github.com/stretchr/testify/require"
)

func TestConverter(t *testing.T) {
	var c mismatched.Converter = &generated.ConverterImpl{}

	input := mismatched.DBCustomers{
		{
			DBPerson: mismatched.DBPerson{
				First: "mary",
				Last:  "brown",
			},
		},
		{
			DBPerson: mismatched.DBPerson{
				First: "john",
				Last:  "smith",
			},
			DBAddress: &mismatched.DBAddress{
				Suburb:   "Surry Hills",
				Postcode: "2010",
			},
		},
	}

	actual := c.Convert(input)

	expected := mismatched.APICustomers{
		{
			APIPerson: &mismatched.APIPerson{
				First: "mary",
				Last:  "brown",
			},
		},
		{
			APIPerson: &mismatched.APIPerson{
				First: "john",
				Last:  "smith",
			},
			APIAddress: mismatched.APIAddress{
				Suburb:   "Surry Hills",
				Postcode: "2010",
			},
		},
	}

	require.Equal(t, expected, actual)
}

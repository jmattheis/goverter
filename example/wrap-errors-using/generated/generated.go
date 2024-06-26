// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package generated

import (
	patherr "github.com/goverter/patherr"
	example "goverter/example"
	"strconv"
)

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source map[int]example.Input) (map[int]example.Output, error) {
	var mapIntExampleOutput map[int]example.Output
	if source != nil {
		mapIntExampleOutput = make(map[int]example.Output, len(source))
		for key, value := range source {
			exampleOutput, err := c.exampleInputToExampleOutput(value)
			if err != nil {
				return nil, patherr.Wrap(err, patherr.Key(key))
			}
			mapIntExampleOutput[key] = exampleOutput
		}
	}
	return mapIntExampleOutput, nil
}
func (c *ConverterImpl) exampleInputToExampleOutput(source example.Input) (example.Output, error) {
	var exampleOutput example.Output
	if source.Friends != nil {
		exampleOutput.Friends = make([]example.Output, len(source.Friends))
		for i := 0; i < len(source.Friends); i++ {
			exampleOutput2, err := c.exampleInputToExampleOutput(source.Friends[i])
			if err != nil {
				return exampleOutput, patherr.Wrap(err, patherr.Field("Friends"), patherr.Index(i))
			}
			exampleOutput.Friends[i] = exampleOutput2
		}
	}
	xint, err := strconv.Atoi(source.Age)
	if err != nil {
		return exampleOutput, patherr.Wrap(err, patherr.Field("Age"))
	}
	exampleOutput.Age = xint
	if source.Attributes != nil {
		exampleOutput.Attributes = make(map[string]int, len(source.Attributes))
		for key, value := range source.Attributes {
			xint2, err := strconv.Atoi(value)
			if err != nil {
				return exampleOutput, patherr.Wrap(err, patherr.Field("Attributes"), patherr.Key(key))
			}
			exampleOutput.Attributes[key] = xint2
		}
	}
	return exampleOutput, nil
}

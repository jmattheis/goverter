// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package generated

import matchtag "github.com/jmattheis/goverter/example/matchTag"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source matchtag.Input) matchtag.Output {
	var exampleOutput matchtag.Output
	exampleOutput.Game = source.Name
	return exampleOutput
}

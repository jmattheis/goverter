// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package generated

import outputraw "github.com/jmattheis/goverter/example/output-raw"

func Hello() string {
	return "World!"
}

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source outputraw.Input) outputraw.Output {
	var rawOutput outputraw.Output
	rawOutput.Name = source.Name
	return rawOutput
}

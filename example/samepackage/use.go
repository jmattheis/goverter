//go:build !goverter

package samepackage

import "errors"

func ValidateAndConvert(source *Input) (*Output, error) {
	if source.Name == "" {
		return nil, errors.New("Name may not be nil")
	}

	c := &ConverterImpl{}
	output := c.Convert(source)
	return output, nil
}

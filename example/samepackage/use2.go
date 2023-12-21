package samepackage

import "errors"

var c Converter

func ValidateAndConvert2(source *Input) (*Output, error) {
	if source.Name == "" {
		return nil, errors.New("Name may not be nil")
	}

	output := c.Convert(source)
	return output, nil
}

// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package generated

import (
	exclude "github.com/jmattheis/goverter/example/enum/exclude"
	"time"
)

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source exclude.MyDuration) time.Duration {
	return time.Duration(source)
}
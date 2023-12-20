//go:build !goverter

package samepackage

func init() {
	c = &ConverterImpl{}
}

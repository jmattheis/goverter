package example

import (
	"bytes"
	"io"
)

// goverter:converter
// goverter:extend ReaderForBytes
type ConverterGlobalMethod interface {
	// goverter:map Content Reader
	Convert(source *Input) *Output
}

func ReaderForBytes(content []byte) io.Reader {
	return bytes.NewReader(content)
}

// You can also specify this method for a specific property conversion

// goverter:converter
type ConverterPropertyMethod interface {
	// goverter:map Content Reader | bytes:NewReader
	Convert(source *Input) *Output
}

type Input struct {
	Source  string
	Content []byte
}

type Output struct {
	Source string
	Reader io.Reader
}

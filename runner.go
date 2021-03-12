package genconv

import (
	"bytes"

	"github.com/jmattheis/go-genconv/internal/comments"
	"github.com/jmattheis/go-genconv/internal/generator"
)

type GenerateConfig struct {
	PackageName string
	ScanDir     string
}

func Generate(c GenerateConfig) ([]byte, error) {
	mapping, err := comments.ParseDocs(c.ScanDir)
	if err != nil {
		return nil, err
	}

	file, err := generator.Generate(c.ScanDir, mapping, generator.Config{Name: c.PackageName})
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = file.Render(buf)
	return buf.Bytes(), err
}

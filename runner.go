package goverter

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmattheis/goverter/comments"
	"github.com/jmattheis/goverter/generator"
)

// GenerateConfig the config for generating a converter.
type GenerateConfig struct {
	PackageName   string
	ScanDir       string
	ExtendMethods []string
}

// GenerateConverter generates converters.
func GenerateConverter(c GenerateConfig) ([]byte, error) {
	mapping, err := comments.ParseDocs(c.ScanDir)
	if err != nil {
		return nil, err
	}

	file, err := generator.Generate(
		c.ScanDir, mapping, generator.Config{Name: c.PackageName, ExtendMethods: c.ExtendMethods})
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = file.Render(buf)
	return buf.Bytes(), err
}

// GenerateConverterFile generates converters and writes them to a file.
func GenerateConverterFile(fileName string, c GenerateConfig) error {
	file, err := GenerateConverter(c)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, file, os.ModePerm)
}

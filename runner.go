package genconv

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmattheis/go-genconv/comments"
	"github.com/jmattheis/go-genconv/generator"
)

type GenerateConfig struct {
	PackageName string
	ScanDir     string
}

func GenerateConverter(c GenerateConfig) ([]byte, error) {
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

func GenerateConverterFile(fileName string, c GenerateConfig) error {
	file, err := GenerateConverter(c)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, file, 0755)
}

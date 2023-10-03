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
	// PackageName is the package to use for the generated code.
	PackageName string
	// ScanDir is the package with golang files to scan for goverter tags.
	ScanDir string
	// ExtendMethods is a list of extensions to load in addition to goverter:extend statements
	// declared on the interface itself.
	ExtendMethods []string
	// WorkingDir is the working directory (usually the location of go.mod file), can be empty.
	WorkingDir string
	// PackagePath is the full package where the generated code is going to be stored, can be empty.
	// PackagePath is needed if goverter:extend tags use exported methods from the package where
	// generated code resides, making sure goverter does not create a loop import statement from
	// the generated code into its own package.
	PackagePath string
	// WrapErrors instructs goverter to wrap conversion errors and return more details when such
	// are available. For examples, for structs, target field name is reported, for slices: index.
	WrapErrors bool
	// IgnoredUnexportedFields tells goverter to ignore fields that are not exported.
	IgnoredUnexportedFields bool
	// MatchFieldsIgnoreCase tells goverter to use case insensitive match everywhere.
	MatchFieldsIgnoreCase bool
	// Add comment on generated structs
	CommentOnStruct string
}

// GenerateConverter generates converters.
func GenerateConverter(c GenerateConfig) ([]byte, error) {
	mapping, err := comments.ParseDocs(comments.ParseDocsConfig{
		PackagePattern: c.ScanDir,
		WorkingDir:     c.WorkingDir,
	})
	if err != nil {
		return nil, err
	}

	file, err := generator.Generate(c.ScanDir, mapping, generator.Config{
		Name:                   c.PackageName,
		PackagePath:            c.PackagePath,
		ExtendMethods:          c.ExtendMethods,
		WorkingDir:             c.WorkingDir,
		WrapErrors:             c.WrapErrors,
		IgnoreUnexportedFields: c.IgnoredUnexportedFields,
		MatchFieldsIgnoreCase:  c.MatchFieldsIgnoreCase,
		CommentOnStruct:        c.CommentOnStruct,
	})
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

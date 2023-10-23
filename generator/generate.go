package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/builder"
	"github.com/jmattheis/goverter/config"
)

// Config the generate config.
type Config struct {
	WorkingDir string
}

// BuildSteps that'll used for generation.
var BuildSteps = []builder.Builder{
	&builder.SkipCopy{},
	&builder.BasicTargetPointerRule{},
	&builder.Pointer{},
	&builder.SourcePointer{},
	&builder.TargetPointer{},
	&builder.Basic{},
	&builder.Struct{},
	&builder.List{},
	&builder.Map{},
}

// Generate generates a jen.File containing converters.
func Generate(converters []*config.Converter, c Config) (map[string][]byte, error) {
	manager := &fileManager{Files: map[string]*managedFile{}}

	for _, converter := range converters {
		jenFile, err := manager.Get(converter, c.WorkingDir)
		if err != nil {
			return nil, err
		}

		if err := generateConverter(converter, c, jenFile); err != nil {
			return nil, err
		}
	}

	return manager.renderFiles()
}

func generateConverter(converter *config.Converter, c Config, f *jen.File) error {
	gen := setupGenerator(converter)

	if err := validateMethods(gen.lookup); err != nil {
		return err
	}

	f.Type().Id(converter.Name).Struct()

	if err := gen.buildMethods(f); err != nil {
		return err
	}
	return nil
}

package generator

import (
	"github.com/jmattheis/goverter/config"
	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/namer"
)

func setupGenerator(converter *config.Converter, n *namer.Namer) (*generator, error) {
	extend := method.NewIndex[method.Definition]()
	for _, def := range converter.Extend {
		extend.RegisterOverrideOverlapping(def, def)
	}

	var err error
	lookup := method.NewIndex[generatedMethod]()
	for _, cMethod := range converter.Methods {
		gen := &generatedMethod{
			Method:   cMethod,
			Dirty:    true,
			Explicit: true,
		}
		gen.IndexID, err = lookup.Register(gen, gen.Definition)
		if err != nil {
			return nil, err
		}
	}

	gen := generator{
		namer:  n,
		conf:   converter,
		lookup: lookup,
		extend: extend,
	}

	return &gen, nil
}

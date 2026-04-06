package generator

import (
	"github.com/jmattheis/goverter/config"
	"github.com/jmattheis/goverter/identity"
	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/namer"
)

func setupGenerator(converter *config.Converter, n *namer.Namer) (*generator, error) {
	extend := method.NewIndex[method.Definition]()
	for _, def := range converter.Extend {
		extend.RegisterOverrideOverlapping(def, def)
	}
	extendIdentity := identity.NewIndex()
	for _, def := range converter.ExtendIdentity {
		extendIdentity.RegisterDefinition(def)
	}

	var err error
	lookup := method.NewIndex[generatedMethod]()
	for _, cMethod := range converter.Methods {
		gen := &generatedMethod{
			Method:   cMethod,
			Dirty:    true,
			Explicit: true,
		}
		if gen.UpdateTarget {
			gen.IndexID, err = lookup.RegisterUpdate(gen, gen.Definition)
		} else {
			gen.IndexID, err = lookup.Register(gen, gen.Definition)
		}
		if err != nil {
			return nil, err
		}
	}

	gen := generator{
		namer:          n,
		conf:           converter,
		lookup:         lookup,
		extend:         extend,
		extendIdentity: extendIdentity,
	}

	return &gen, nil
}

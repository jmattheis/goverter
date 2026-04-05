package identity

import (
	"fmt"
	"go/types"

	"github.com/jmattheis/goverter/xtype"
)

type ParseOpts struct {
	Location string
	// Converter         types.Type
	OutputPackagePath string

	ErrorPrefix string
	// Params            ParamType
	// ParamsMultiSource bool
	// AllowTypeParams   bool
	//
	// ContextMatch *regexp.Regexp
	//
	// Generated   bool
	// CustomCall  *jen.Statement
	// UpdateParam string
}

func Parse(obj types.Object, opts *ParseOpts) (*Definition, error) {
	identityDef := &Definition{
		ID:       obj.String(),
		OriginID: obj.String(),
		Name:     obj.Name(),
	}

	formatErr := func(s string) error {
		loc := ""
		if opts.Location != "" {
			loc = opts.Location + "\n    "
		}
		return fmt.Errorf("%s:\n    %s%s%s\n\n%s", opts.ErrorPrefix, loc, obj.String(), identityDef.ArgDebug("        "), s)
	}

	if !xtype.Accessible(obj, opts.OutputPackagePath) {
		return nil, formatErr("must be exported")
	}

	// sig, ok := obj.Type().(*types.Signature)
	// if !ok {
	// 	return nil, formatErr("must be a function")
	// }

	if pkg := obj.Pkg(); pkg != nil {
		identityDef.Package = pkg.Path()
	}

	identityDef.Type = xtype.TypeOf(obj.Type())

	return identityDef, nil
}

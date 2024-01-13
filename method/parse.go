package method

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

type ParamType int

const (
	ParamsRequired ParamType = iota
	ParamsOptional
	ParamsNone
)

type ParseOpts struct {
	Location          string
	Converter         types.Type
	OutputPackagePath string

	ErrorPrefix     string
	Params          ParamType
	AllowTypeParams bool

	Generated  bool
	CustomCall *jen.Statement
}

// Parse parses an function into a Definition.
func Parse(obj types.Object, opts *ParseOpts) (*Definition, error) {
	formatErr := func(s string) error {
		loc := ""
		if opts.Location != "" {
			loc = opts.Location + "\n    "
		}
		return fmt.Errorf("%s:\n    %s%s\n\n%s", opts.ErrorPrefix, loc, obj.String(), s)
	}

	if !xtype.Accessible(obj, opts.OutputPackagePath) {
		return nil, formatErr("must be exported")
	}

	sig, ok := obj.Type().(*types.Signature)
	if !ok {
		return nil, formatErr("must be a function")
	}
	if sig.Results().Len() == 0 || sig.Results().Len() > 2 {
		return nil, formatErr("must have one or two returns")
	}
	returnError := false
	if sig.Results().Len() == 2 {
		if i, ok := sig.Results().At(1).Type().(*types.Named); ok && i.Obj().Name() == "error" && i.Obj().Pkg() == nil {
			returnError = true
		} else {
			return nil, formatErr("must have type error as second return but has: " + sig.Results().At(1).Type().String())
		}
	}

	methodDef := &Definition{
		ID:         obj.String(),
		OriginID:   obj.String(),
		Generated:  opts.Generated,
		CustomCall: opts.CustomCall,
		Parameters: Parameters{
			ReturnError: returnError,
			Target:      xtype.TypeOf(sig.Results().At(0).Type()),
			TypeParams:  sig.TypeParams().Len() > 0,
		},
		Name: obj.Name(),
	}

	if methodDef.TypeParams && !opts.AllowTypeParams {
		return nil, formatErr("must not be generic")
	}

	if pkg := obj.Pkg(); pkg != nil {
		methodDef.Package = pkg.Path()
	}

	if opts.Params == ParamsNone && sig.Params().Len() > 0 {
		return nil, formatErr("must have no parameters")
	}

	switch sig.Params().Len() {
	case 2:
		if opts.Converter == nil {
			return nil, formatErr("must have one parameter")
		}

		actual := sig.Params().At(0).Type().String()
		if actual != opts.Converter.String() {
			return nil, formatErr(
				fmt.Sprintf("first parameter must be of type %s but was %s when having two parameters", opts.Converter.String(), actual))
		}
		methodDef.Parameters.SelfAsFirstParameter = true
		methodDef.Parameters.Source = xtype.TypeOf(sig.Params().At(1).Type())
	case 1:
		methodDef.Parameters.Source = xtype.TypeOf(sig.Params().At(0).Type())
	case 0:
		if opts.Params == ParamsRequired {
			return nil, formatErr("must have at least one parameter")
		}
	default:
		return nil, formatErr("must have one or two parameters")
	}

	return methodDef, nil
}

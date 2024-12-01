package method

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"

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

	ErrorPrefix       string
	Params            ParamType
	ParamsMultiSource bool
	AllowTypeParams   bool

	ContextMatch *regexp.Regexp

	Generated   bool
	CustomCall  *jen.Statement
	UpdateParam string
}

// Parse parses an function into a Definition.
func Parse(obj types.Object, opts *ParseOpts) (*Definition, error) {
	methodDef := &Definition{
		ID:         obj.String(),
		OriginID:   obj.String(),
		Generated:  opts.Generated,
		CustomCall: opts.CustomCall,
		Parameters: Parameters{
			Context: make(map[string]*xtype.Type, 0),
		},
		Name: obj.Name(),
	}

	formatErr := func(s string) error {
		loc := ""
		if opts.Location != "" {
			loc = opts.Location + "\n    "
		}
		return fmt.Errorf("%s:\n    %s%s%s\n\n%s", opts.ErrorPrefix, loc, obj.String(), methodDef.ArgDebug("        "), s)
	}

	if !xtype.Accessible(obj, opts.OutputPackagePath) {
		return nil, formatErr("must be exported")
	}

	sig, ok := obj.Type().(*types.Signature)
	if !ok {
		return nil, formatErr("must be a function")
	}
	resultsLen := sig.Results().Len()

	methodDef.TypeParams = sig.TypeParams().Len() > 0

	if pkg := obj.Pkg(); pkg != nil {
		methodDef.Package = pkg.Path()
	}

	for i := 0; i < sig.Params().Len(); i++ {
		arg := Arg{
			Name: sig.Params().At(i).Name(),
			Type: xtype.TypeOf(sig.Params().At(i).Type()),
		}

		switch {
		case types.Identical(arg.Type.T, opts.Converter):
			arg.Use = ArgUseInterface
		case opts.UpdateParam != "" && arg.Name == opts.UpdateParam:
			arg.Use = ArgUseTarget
			methodDef.Target = arg.Type
			methodDef.UpdateTarget = true

			switch {
			case resultsLen == 0:
				// okay nothing more
			case resultsLen == 1 && isError(sig.Results().At(0)):
				methodDef.ReturnError = true
			default:
				return nil, formatErr("The signature one non 'error' result or multiple results is not supported for goverter:update signatures.")
			}
		case opts.ContextMatch.MatchString(arg.Name):
			methodDef.Context[arg.Type.String] = arg.Type
			arg.Use = ArgUseContext
		case methodDef.Source == nil:
			arg.Use = ArgUseSource
			methodDef.Source = arg.Type
			methodDef.Signature.Source = methodDef.Source.String
		default:
			arg.Use = ArgUseMultiSource
			methodDef.MultiSources = append(methodDef.MultiSources, arg.Type)
		}

		methodDef.RawArgs = append(methodDef.RawArgs, arg)
	}
	if !methodDef.UpdateTarget && opts.UpdateParam != "" {
		return nil, formatErr(fmt.Sprintf("Argument %q must exist when using 'goverter:target %s'", opts.UpdateParam, opts.UpdateParam))
	}

	if !methodDef.UpdateTarget {
		if resultsLen == 0 || resultsLen > 2 {
			return nil, formatErr("must have one or two returns")
		}
		if resultsLen == 2 {
			if isError(sig.Results().At(1)) {
				methodDef.ReturnError = true
			} else {
				return nil, formatErr("must have type error as second return but has: " + sig.Results().At(1).Type().String())
			}
		}

		methodDef.Target = xtype.TypeOf(sig.Results().At(0).Type())
	}

	if methodDef.TypeParams && !opts.AllowTypeParams {
		return nil, formatErr("must not be generic")
	}

	switch {
	case opts.Params == ParamsNone && methodDef.Source != nil:
		return nil, formatErr("must have no source params")
	case opts.Params == ParamsRequired && methodDef.Source == nil:
		return nil, formatErr("must have at least one source param")
	case !opts.ParamsMultiSource && len(methodDef.MultiSources) > 0:
		return nil, formatErr("must have only one source param")
	}

	methodDef.Signature.Target = methodDef.Target.String

	return methodDef, nil
}

func isError(obj *types.Var) bool {
	t, ok := obj.Type().(*types.Named)
	return ok && t.Obj().Name() == "error" && t.Obj().Pkg() == nil
}

func (def *Definition) ArgDebug(indent string) string {
	var lines []string
	for _, arg := range def.RawArgs {
		argUse := arg.Use
		if arg.Use == ArgUseMultiSource {
			argUse = ArgUseSource
		} else if arg.Use == ArgUseInterface {
			argUse = ArgUseContext
		}
		lines = append(lines, fmt.Sprintf("[%s] %s", argUse, arg.Type.String))
	}

	if def.Target != nil && !def.UpdateTarget {
		lines = append(lines, fmt.Sprintf("[target] %s", def.Target.String))
	}

	if len(lines) == 0 {
		return ""
	}

	return "\n" + indent + strings.Join(lines, "\n"+indent)
}

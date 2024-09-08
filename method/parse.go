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

	Generated  bool
	CustomCall *jen.Statement
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
	if sig.Results().Len() == 0 || sig.Results().Len() > 2 {
		return nil, formatErr("must have one or two returns")
	}
	if sig.Results().Len() == 2 {
		if i, ok := sig.Results().At(1).Type().(*types.Named); ok && i.Obj().Name() == "error" && i.Obj().Pkg() == nil {
			methodDef.ReturnError = true
		} else {
			return nil, formatErr("must have type error as second return but has: " + sig.Results().At(1).Type().String())
		}
	}

	methodDef.Target = xtype.TypeOf(sig.Results().At(0).Type())
	methodDef.TypeParams = sig.TypeParams().Len() > 0

	if pkg := obj.Pkg(); pkg != nil {
		methodDef.Package = pkg.Path()
	}

	for i := 0; i < sig.Params().Len(); i++ {
		arg := Arg{
			Name: sig.Params().At(i).Name(),
			Type: xtype.TypeOf(sig.Params().At(i).Type()),
		}

		if types.Identical(arg.Type.T, opts.Converter) {
			arg.Use = ArgUseInterface
		} else if opts.ContextMatch.MatchString(arg.Name) {
			methodDef.Context[arg.Type.String] = arg.Type
			arg.Use = ArgUseContext
		} else if methodDef.Source == nil {
			arg.Use = ArgUseSource
			methodDef.Source = arg.Type
			methodDef.Signature.Source = methodDef.Source.String
		} else {
			arg.Use = ArgUseMultiSource
			methodDef.MultiSources = append(methodDef.MultiSources, arg.Type)
		}

		methodDef.RawArgs = append(methodDef.RawArgs, arg)
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

	if def.Target != nil {
		lines = append(lines, fmt.Sprintf("[target] %s", def.Target.String))
	}

	if len(lines) == 0 {
		return ""
	}

	return "\n" + indent + strings.Join(lines, "\n"+indent)
}

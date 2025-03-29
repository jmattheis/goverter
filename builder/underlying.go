package builder

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// UseUnderlyingTypeMethods handles UseUnderlyingTypeMethods.
type UseUnderlyingTypeMethods struct{}

// Matches returns true, if the builder can create handle the given types.
func (*UseUnderlyingTypeMethods) Matches(ctx *MethodContext, source, target *xtype.Type) bool {
	if !ctx.Conf.UseUnderlyingTypeMethods {
		return false
	}

	sourceUnderlying, targetUnderlying := findUnderlyingExtendMapping(ctx, source, target)
	return sourceUnderlying || targetUnderlying
}

// Build creates conversion source code for the given source and target type.
func (*UseUnderlyingTypeMethods) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	if isEnum(ctx, source, target) {
		return nil, nil, NewError(fmt.Sprintf(`The conversion between the types
    %s
    %s

does qualify for enum conversion but also match an extend method via useUnderlyingTypeMethods.
You have to disable enum or useUnderlyingTypeMethods to resolve the setting conflict.`, source.String, target.String))
	}

	sourceUnderlying, targetUnderlying := findUnderlyingExtendMapping(ctx, source, target)

	innerSource := source
	innerTarget := target

	if sourceUnderlying {
		innerSource = xtype.TypeOf(source.NamedType.Underlying())
		sourceID = xtype.OtherID(innerSource.TypeAsJen().Call(sourceID.Code))
	}

	if targetUnderlying {
		innerTarget = xtype.TypeOf(target.NamedType.Underlying())
	}

	stmt, id, err := gen.Build(ctx, sourceID, innerSource, innerTarget, errPath)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: innerSource.String,
			TargetID:   "*",
			TargetType: innerTarget.String,
		})
	}

	if targetUnderlying {
		id = xtype.OtherID(target.TypeAsJen().Call(id.Code))
	}

	return stmt, id, err
}

func (u *UseUnderlyingTypeMethods) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *Error) {
	return AssignByBuild(u, gen, ctx, assignTo, sourceID, source, target, errPath)
}

func findUnderlyingExtendMapping(ctx *MethodContext, source, target *xtype.Type) (underlyingSource, underlyingTarget bool) {
	if source.Named {
		if ctx.HasMethod(ctx, source.NamedType.Underlying(), target.T) {
			return true, false
		}

		if target.Named && ctx.HasMethod(ctx, source.NamedType.Underlying(), target.NamedType.Underlying()) {
			return true, true
		}
	}

	if target.Named && ctx.HasMethod(ctx, source.T, target.NamedType.Underlying()) {
		return false, true
	}

	return false, false
}

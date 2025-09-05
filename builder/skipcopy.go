package builder

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// SkipCopy handles FlagSkipCopySameType.
type SkipCopy struct{}

// Matches returns true, if the builder can create handle the given types.
func (*SkipCopy) Matches(ctx *MethodContext, source, target *xtype.Type) bool {
	return ctx.Conf.SkipCopySameType && types.Identical(source.T, target.T)
}

// Build creates conversion source code for the given source and target type.
func (*SkipCopy) Build(_ Generator, _ *MethodContext, sourceID *xtype.JenID, _, _ *xtype.Type, _ ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	return nil, sourceID, nil
}

func (*SkipCopy) Assign(_ Generator, _ *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, _, _ *xtype.Type, _ ErrorPath) ([]jen.Code, *Error) {
	return []jen.Code{assignTo.Stmt.Clone().Op("=").Add(sourceID.Code)}, nil
}

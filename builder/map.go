package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// Map handles map types.
type Map struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Map) Matches(source, target *xtype.Type) bool {
	return source.Map && target.Map
}

// Build creates conversion source code for the given source and target type.
func (*Map) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	targetMap := ctx.Name(target.ID())
	key, value := ctx.Map()

	block, newKey, err := gen.Build(ctx, xtype.VariableID(jen.Id(key)), source.MapKey, target.MapKey, NoWrap)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: "<mapkey> " + source.MapKey.T.String(),
			TargetID:   "[]",
			TargetType: "<mapkey> " + target.MapKey.T.String(),
		})
	}
	valueStmt, valueKey, err := gen.Build(
		ctx, xtype.VariableID(jen.Id(value)), source.MapValue, target.MapValue, NoWrap)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: "<mapvalue> " + source.MapValue.T.String(),
			TargetID:   "[]",
			TargetType: "<mapvalue> " + target.MapValue.T.String(),
		})
	}
	block = append(block, valueStmt...)
	block = append(block, jen.Id(targetMap).Index(newKey.Code).Op("=").Add(valueKey.Code))

	stmt := []jen.Code{
		jen.Id(targetMap).Op(":=").Make(target.TypeAsJen(), jen.Len(sourceID.Code.Clone())),
		jen.For(jen.List(jen.Id(key), jen.Id(value)).Op(":=").Range().Add(sourceID.Code)).
			Block(block...),
	}

	return stmt, xtype.VariableID(jen.Id(targetMap)), nil
}

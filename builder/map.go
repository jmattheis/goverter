package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// Map handles map types.
type Map struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Map) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	return source.Map && target.Map
}

// Build creates conversion source code for the given source and target type.
func (m *Map) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	return BuildByAssign(m, gen, ctx, sourceID, source, target, errPath)
}

func (*Map) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *Error) {
	ctx.SetErrorTargetVar(jen.Nil())
	key, value := ctx.Map()

	errPath = errPath.Key(jen.Id(key))

	block, keyID, err := gen.Build(ctx, xtype.VariableID(jen.Id(key)), source.MapKey, target.MapKey, errPath)
	if err != nil {
		return nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: "<mapkey> " + source.MapKey.String,
			TargetID:   "[]",
			TargetType: "<mapkey> " + target.MapKey.String,
		})
	}
	valueStmt, err := gen.Assign(
		ctx, assignTo.WithIndex(keyID.Code).WithMust(), xtype.VariableID(jen.Id(value)), source.MapValue, target.MapValue, errPath)
	if err != nil {
		return nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: "<mapvalue> " + source.MapValue.String,
			TargetID:   "[]",
			TargetType: "<mapvalue> " + target.MapValue.String,
		})
	}
	block = append(block, valueStmt...)

	stmt := []jen.Code{
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(
			assignTo.Stmt.Clone().Op("=").Make(target.TypeAsJen(), jen.Len(sourceID.Code.Clone())),
			jen.For(jen.List(jen.Id(key), jen.Id(value)).Op(":=").Range().Add(sourceID.Code)).
				Block(block...),
		),
	}

	return stmt, nil
}

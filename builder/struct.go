package builder

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

type Struct struct{}

func (*Struct) Matches(source, target *Type) bool {
	return source.Struct && target.Struct
}

func (*Struct) Build(gen Generator, ctx *MethodContext, sourceID *JenID, source, target *Type) ([]jen.Code, *JenID, *Error) {
	name := ctx.Of(target, "structTarget")
	stmt := []jen.Code{
		jen.Var().Id(name).Add(target.TypeAsJen()),
	}

	sourceStruct := source.StructType
	targetStruct := target.StructType

	sourceMethods := map[string]*types.Var{}
	for i := 0; i < sourceStruct.NumFields(); i++ {
		m := sourceStruct.Field(i)
		sourceMethods[m.Name()] = m
	}

	for i := 0; i < targetStruct.NumFields(); i++ {
		targetField := targetStruct.Field(i)
		sourceField, ok := sourceMethods[targetField.Name()]
		if !ok {
			cause := fmt.Sprintf("Cannot set value for field %s because no it does not exist on the source entry", targetField.Name())
			return nil, nil, NewError(cause).Lift(&Path{
				Prefix:     ".",
				SourceID:   "???",
				TargetID:   targetField.Name(),
				TargetType: targetField.Type().String(),
			})
		}

		fieldSourceId := sourceID.Code.Clone().Dot(sourceField.Name())

		fieldStmt, fieldID, err := gen.Build(ctx, VariableID(fieldSourceId), TypeOf(sourceField.Type()), TypeOf(targetField.Type()))
		if err != nil {
			return nil, nil, err.Lift(&Path{
				Prefix:     ".",
				SourceID:   sourceField.Name(),
				TargetID:   targetField.Name(),
				TargetType: targetField.Type().String(),
				SourceType: sourceField.Type().String(),
			})
		}
		stmt = append(stmt, fieldStmt...)
		stmt = append(stmt, jen.Id(name).Dot(targetField.Name()).Op("=").Add(fieldID.Code))
	}

	return stmt, VariableID(jen.Id(name)), nil
}

package builder

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/xtype"
)

// Struct handles struct types.
type Struct struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Struct) Matches(source, target *xtype.Type) bool {
	return source.Struct && target.Struct
}

// Build creates conversion source code for the given source and target type.
func (*Struct) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	name := ctx.Name(target.ID())
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
		if _, ignore := ctx.IgnoredFields[targetField.Name()]; ignore {
			continue
		}
		if !targetField.Exported() {
			cause := unexportedStructError(targetField.Name(), source.T.String(), target.T.String())
			return nil, nil, NewError(cause).Lift(&Path{
				Prefix:     ".",
				SourceID:   "???",
				TargetID:   targetField.Name(),
				TargetType: targetField.Type().String(),
			})
		}

		sourceName := targetField.Name()
		if ctx.Signature.Target == target.T.String() {
			if override, ok := ctx.Mapping[targetField.Name()]; ok {
				sourceName = override
			}
		}
		sourceField, ok := sourceMethods[sourceName]
		if !ok {
			cause := fmt.Sprintf("Cannot set value for field %s because it does not exist on the source entry", targetField.Name())
			return nil, nil, NewError(cause).Lift(&Path{
				Prefix:     ".",
				SourceID:   "???",
				TargetID:   targetField.Name(),
				TargetType: targetField.Type().String(),
			})
		}

		fieldSourceID := sourceID.Code.Clone().Dot(sourceField.Name())

		fieldStmt, fieldID, err := gen.Build(ctx, xtype.VariableID(fieldSourceID), xtype.TypeOf(sourceField.Type()), xtype.TypeOf(targetField.Type()))
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

	return stmt, xtype.VariableID(jen.Id(name)), nil
}

func unexportedStructError(targetField, sourceType, targetType string) string {
	return fmt.Sprintf(`Cannot set value for unexported field "%s".

Possible solutions:

* Ignore the given field with:

      // goverter:ignore %s

* Convert the struct yourself and use goverter for converting nested structs / maps / lists.

* Create a custom converter function (only works, if the struct with unexported fields is nested inside another struct)

      func CustomConvert(source %s) %s {
          // implement me
      }

      // goverter:extend CustomConvert
      type MyConverter interface {
          // ...
      }

See https://github.com/jmattheis/goverter#extend-with-custom-implementation`, targetField, targetField, sourceType, targetType)
}

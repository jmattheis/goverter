package builder

import (
	"fmt"
	"go/types"
	"strings"

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

	for i := 0; i < target.StructType.NumFields(); i++ {
		targetField := target.StructType.Field(i)
		if _, ignore := ctx.IgnoredFields[targetField.Name()]; ignore {
			continue
		}
		if !targetField.Exported() && ctx.IgnoreUnexportedFields {
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

		targetFieldType := xtype.TypeOf(targetField.Type())

		if def, ok := ctx.ExtendMapping[targetField.Name()]; ok {
			params := []jen.Code{}
			if def.SelfAsFirstParam {
				params = append(params, jen.Id(xtype.ThisVar))
			}

			if def.Source != nil {
				if def.Source.T.String() != source.T.String() {
					cause := fmt.Sprintf("cannot not use\n\t%s\nbecause source type mismatch\n\nExtend method param type: %s\nConverter source type: %s", def.ID, def.Source.T.String(), source.T.String())
					return nil, nil, NewError(cause).Lift(&Path{
						Prefix:     ".",
						SourceID:   "<mapExtend>",
						SourceType: def.ID,
						TargetID:   targetField.Name(),
						TargetType: targetField.Type().String(),
					})
				}
				params = append(params, sourceID.Code.Clone())
			}

			if def.Target.T.String() != targetFieldType.T.String() {
				cause := fmt.Sprintf("Extend method return type mismatches with target: %s != %s", def.Target.T.String(), targetFieldType.T.String())
				return nil, nil, NewError(cause).Lift(&Path{
					Prefix:     ".",
					SourceID:   "()",
					SourceType: def.Target.T.String(),
					TargetID:   targetField.Name(),
					TargetType: targetField.Type().String(),
				}).Lift(&Path{
					Prefix:     ".",
					SourceID:   "<mapExtend>",
					SourceType: def.ID,
				})
			}
			stmt = append(stmt, jen.Id(name).Dot(targetField.Name()).Op("=").Add(def.Call.Clone().Call(params...)))
			continue
		}

		// To find out the source code for an error message like "error setting field PostalCode", people sometimes
		// search their codebase/repo for the exact message match, in full. Inline the error message AS IS into the
		// generated code to satisfy this search and speed up the troubleshooting efforts.
		errWrapper := Wrap("error setting field " + targetField.Name())
		if _, ok := ctx.IdentityMapping[targetField.Name()]; ok {
			fieldStmt, fieldID, err := gen.Build(ctx, sourceID, source, targetFieldType, errWrapper)
			if err != nil {
				return nil, nil, err.Lift(&Path{
					Prefix:     ".",
					SourceID:   "<mapIdentity>",
					SourceType: source.T.String(),
					TargetID:   targetField.Name(),
					TargetType: targetField.Type().String(),
				})
			}
			stmt = append(stmt, fieldStmt...)
			stmt = append(stmt, jen.Id(name).Dot(targetField.Name()).Op("=").Add(fieldID.Code))
			continue
		}

		nextID, nextSource, mapStmt, lift, err := mapField(gen, ctx, targetField, sourceID, source, target)
		if err != nil {
			return nil, nil, err
		}
		stmt = append(stmt, mapStmt...)

		fieldStmt, fieldID, err := gen.Build(ctx, xtype.VariableID(nextID), nextSource, targetFieldType, errWrapper)
		if err != nil {
			return nil, nil, err.Lift(lift...)
		}
		stmt = append(stmt, fieldStmt...)
		stmt = append(stmt, jen.Id(name).Dot(targetField.Name()).Op("=").Add(fieldID.Code))
	}

	return stmt, xtype.VariableID(jen.Id(name)), nil
}

func mapField(gen Generator, ctx *MethodContext, targetField *types.Var, sourceID *xtype.JenID, source, target *xtype.Type) (*jen.Statement, *xtype.Type, []jen.Code, []*Path, *Error) {
	lift := []*Path{}

	mappedName, hasOverride := ctx.Mapping[targetField.Name()]
	if ctx.Signature.Target != target.T.String() || !hasOverride {
		sourceMatch, err := source.StructField(targetField.Name(), ctx.MatchIgnoreCase, ctx.IgnoredFields)
		if err == nil {
			nextID := sourceID.Code.Clone().Dot(sourceMatch.Name)
			lift = append(lift, &Path{
				Prefix:     ".",
				SourceID:   sourceMatch.Name,
				SourceType: sourceMatch.Type.T.String(),
				TargetID:   targetField.Name(),
				TargetType: targetField.Type().String(),
			})
			return nextID, sourceMatch.Type, []jen.Code{}, lift, nil
		}
		// field lookup either did not find anything or failed due to ambiguous match with case ignored
		cause := fmt.Sprintf("Cannot match the target field with the source entry: %s.", err.Error())
		return nil, nil, nil, nil, NewError(cause).Lift(&Path{
			Prefix:     ".",
			SourceID:   "???",
			TargetID:   targetField.Name(),
			TargetType: targetField.Type().String(),
		})
	}

	path := strings.Split(mappedName, ".")
	var condition *jen.Statement

	stmt := []jen.Code{}
	nextID := sourceID.Code
	nextSource := source
	for i := 0; i < len(path); i++ {
		if nextSource.Pointer {
			addCondition := nextID.Clone().Op("!=").Nil()
			if condition == nil {
				condition = addCondition
			} else {
				condition = condition.Clone().Op("&&").Add(addCondition)
			}
			nextSource = nextSource.PointerInner
		}
		if !nextSource.Struct {
			cause := fmt.Sprintf("Cannot access '%s' on %s.", path[i], nextSource.T)
			return nil, nil, nil, nil, NewError(cause).Lift(&Path{
				Prefix:     ".",
				SourceID:   path[i],
				SourceType: "???",
			}).Lift(lift...)
		}
		// since we are searching for a mapped name, search for exact match, explicit field map does not ignore case
		sourceMatch, err := nextSource.StructField(path[i], false, ctx.IgnoredFields)
		if err == nil {
			nextSource = sourceMatch.Type
			nextID = nextID.Clone().Dot(sourceMatch.Name)
			liftPath := &Path{
				Prefix:     ".",
				SourceID:   sourceMatch.Name,
				SourceType: nextSource.T.String(),
			}

			if i == len(path)-1 {
				liftPath.TargetID = targetField.Name()
				liftPath.TargetType = targetField.Type().String()
				if condition != nil && !nextSource.Pointer {
					liftPath.SourceType = fmt.Sprintf("*%s (It is a pointer because the nested property in the goverter:map was a pointer)", liftPath.SourceType)
				}
			}
			lift = append(lift, liftPath)
			continue
		}

		cause := fmt.Sprintf("Cannot find the mapped field on the source entry: %s.", err.Error())
		return nil, nil, []jen.Code{}, nil, NewError(cause).Lift(&Path{
			Prefix:     ".",
			SourceID:   path[i],
			SourceType: "???",
		}).Lift(lift...)
	}
	if condition != nil {
		pointerNext := nextSource
		if !nextSource.Pointer {
			pointerNext = xtype.TypeOf(types.NewPointer(nextSource.T))
		}
		tempName := ctx.Name(pointerNext.ID())
		stmt = append(stmt, jen.Var().Id(tempName).Add(pointerNext.TypeAsJen()))
		if nextSource.Pointer {
			stmt = append(stmt, jen.If(condition).Block(
				jen.Id(tempName).Op("=").Add(nextID.Clone()),
			))
		} else {
			stmt = append(stmt, jen.If(condition).Block(
				jen.Id(tempName).Op("=").Op("&").Add(nextID.Clone()),
			))
		}
		nextSource = pointerNext
		nextID = jen.Id(tempName)
	}

	return nextID, nextSource, stmt, lift, nil
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

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
func (*Struct) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	return source.Struct && target.Struct
}

// Build creates conversion source code for the given source and target type.
func (*Struct) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	name := ctx.Name(target.ID())
	ctx.SetErrorTargetVar(jen.Id(name))
	stmt := []jen.Code{
		jen.Var().Id(name).Add(target.TypeAsJen()),
	}

	additionalFieldSources, err := parseAutoMap(ctx, source)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < target.StructType.NumFields(); i++ {
		targetField := target.StructType.Field(i)

		fieldMapping := ctx.Field(target, targetField.Name())

		if fieldMapping.Ignore {
			continue
		}
		if !targetField.Exported() && ctx.Flags.Has(FlagIgnoreUnexported) {
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

		// To find out the source code for an error message like "error setting field PostalCode", people sometimes
		// search their codebase/repo for the exact message match, in full. Inline the error message AS IS into the
		// generated code to satisfy this search and speed up the troubleshooting efforts.
		errWrapper := Wrap("error setting field " + targetField.Name())

		if fieldMapping.Function == nil {
			nextID, nextSource, mapStmt, lift, skip, err := mapField(gen, ctx, targetField, sourceID, source, target, additionalFieldSources)
			if skip {
				continue
			}
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
		} else {
			def := fieldMapping.Function

			sourceLift := []*Path{}
			var functionCallSourceID *xtype.JenID
			var functionCallSourceType *xtype.Type
			if def.Source != nil {
				nextID, nextSource, mapStmt, mapLift, _, err := mapField(gen, ctx, targetField, sourceID, source, target, additionalFieldSources)
				if err != nil {
					return nil, nil, err
				}
				sourceLift = mapLift
				stmt = append(stmt, mapStmt...)

				switch {
				case def.Source.T.String() == nextSource.T.String():
					functionCallSourceID = xtype.VariableID(nextID.Clone())
					functionCallSourceType = nextSource
				case fieldMapping.Source == "." && sourceID.ParentPointer != nil &&
					def.Source.T.String() == source.AsPointer().T.String():
					functionCallSourceID = sourceID.ParentPointer
					functionCallSourceType = source.AsPointer()
				default:
					cause := fmt.Sprintf("cannot not use\n\t%s\nbecause source type mismatch\n\nExtend method param type: %s\nConverter source type: %s", def.ID, def.Source.T.String(), nextSource.T.String())
					return nil, nil, NewError(cause).Lift(&Path{
						Prefix:     "(",
						SourceID:   "source)",
						SourceType: def.Source.T.String(),
					}).Lift(&Path{
						Prefix:     ":",
						SourceID:   def.Name,
						SourceType: def.ID,
					}).Lift(sourceLift...)
				}
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
					Prefix:     ":",
					SourceID:   def.Name,
					SourceType: def.ID,
				}).Lift(sourceLift...)
			}
			callStmt, callReturnID, err := gen.CallExtendMethod(ctx, fieldMapping.Function, functionCallSourceID, functionCallSourceType, targetFieldType, errWrapper)
			if err != nil {
				return nil, nil, err.Lift(sourceLift...)
			}
			stmt = append(stmt, callStmt...)
			stmt = append(stmt, jen.Id(name).Dot(targetField.Name()).Op("=").Add(callReturnID.Code))
		}
	}

	return stmt, xtype.VariableID(jen.Id(name)), nil
}

func mapField(
	gen Generator,
	ctx *MethodContext,
	targetField *types.Var,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	additionalFieldSources []xtype.FieldSources,
) (*jen.Statement, *xtype.Type, []jen.Code, []*Path, bool, *Error) {
	lift := []*Path{}
	ignored := func(name string) bool {
		return ctx.Field(target, name).Ignore
	}

	def := ctx.Field(target, targetField.Name())
	pathString := def.Source
	if pathString == "." {
		lift = append(lift, &Path{
			Prefix:     ".",
			SourceID:   " ",
			SourceType: "goverter:map . " + targetField.Name(),
			TargetID:   targetField.Name(),
			TargetType: targetField.Type().String(),
		})
		return sourceID.Code, source, nil, lift, false, nil
	}

	var path []string
	if pathString == "" {
		sourceMatch, err := xtype.FindField(targetField.Name(), ctx.Flags.Has(FlagMatchIgnoreCase), ignored, source, additionalFieldSources)
		if err != nil {
			cause := fmt.Sprintf("Cannot match the target field with the source entry: %s.", err.Error())
			skip := false
			if ctx.Flags.Has(FlagIgnoreMissing) {
				_, skip = err.(*xtype.NoMatchError)
			}
			return nil, nil, nil, nil, skip, NewError(cause).Lift(&Path{
				Prefix:     ".",
				SourceID:   "???",
				TargetID:   targetField.Name(),
				TargetType: targetField.Type().String(),
			})
		}

		path = sourceMatch.Path
	} else {
		path = strings.Split(pathString, ".")
	}

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
			return nil, nil, nil, nil, false, NewError(cause).Lift(&Path{
				Prefix:     ".",
				SourceID:   path[i],
				SourceType: "???",
			}).Lift(lift...)
		}
		sourceMatch, err := xtype.FindExactField(nextSource, path[i])
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
		return nil, nil, []jen.Code{}, nil, false, NewError(cause).Lift(&Path{
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

	return nextID, nextSource, stmt, lift, false, nil
}

func parseAutoMap(ctx *MethodContext, source *xtype.Type) ([]xtype.FieldSources, *Error) {
	fieldSources := []xtype.FieldSources{}
	for _, field := range ctx.AutoMap {
		innerSource := source
		lift := []*Path{}
		path := strings.Split(field, ".")
		for _, part := range path {
			field, err := xtype.FindExactField(innerSource, part)
			if err != nil {
				return nil, NewError(err.Error()).Lift(&Path{
					Prefix:     ".",
					SourceID:   part,
					SourceType: "goverter:autoMap",
				}).Lift(lift...)
			}
			lift = append(lift, &Path{
				Prefix:     ".",
				SourceID:   field.Name,
				SourceType: field.Type.T.String(),
			})
			innerSource = field.Type

			switch {
			case innerSource.Pointer && innerSource.PointerInner.Struct:
				innerSource = xtype.TypeOf(innerSource.PointerInner.StructType)
			case innerSource.Struct:
				// ok
			default:
				return nil, NewError(fmt.Sprintf("%s is not a struct or struct pointer", part)).Lift(lift...)
			}
		}

		fieldSources = append(fieldSources, xtype.FieldSources{Path: path, Type: innerSource})
	}
	return fieldSources, nil
}

func unexportedStructError(targetField, sourceType, targetType string) string {
	return fmt.Sprintf(`Cannot set value for unexported field "%s".

Possible solutions:

* Ignore the given field:
  https://goverter.jmattheis.de/#/conversion/mapping?id=ignore

* Create a custom converter function:
  https://goverter.jmattheis.de/#/conversion/custom`, targetField)
}

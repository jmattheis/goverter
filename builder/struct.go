package builder

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/goverter/config"
	"github.com/jmattheis/goverter/method"
	"github.com/jmattheis/goverter/xtype"
)

// Struct handles struct types.
type Struct struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Struct) Matches(_ *MethodContext, source, target *xtype.Type) bool {
	return source.Struct && target.Struct
}

// Build creates conversion source code for the given source and target type.
func (s *Struct) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *xtype.JenID, *Error) {
	// Optimization for golang sets
	if !source.Named && !target.Named && source.StructType.NumFields() == 0 && target.StructType.NumFields() == 0 {
		return nil, sourceID, nil
	}
	return BuildByAssign(s, gen, ctx, sourceID, source, target, errPath)
}

func (s *Struct) Assign(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath) ([]jen.Code, *Error) {
	additionalFieldSources, err := parseAutoMap(ctx, source)
	if err != nil {
		return nil, err
	}

	stmt := []jen.Code{}

	definedFields := ctx.DefinedFields(target)
	// nil indicates that the field was found in one of the paths
	// non nil indicates that it was not found and should error if it's not found in the alternative path
	missing := map[string]*Error{}

	fieldGens := []func() ([]jen.Code, bool, *Error){
		func() ([]jen.Code, bool, *Error) {
			return fields(gen, ctx, assignTo, sourceID, source, target, errPath, additionalFieldSources, definedFields, missing)
		},
	}

	if ctx.Conf.SettersEnabled {
		fieldGens = append(fieldGens, func() ([]jen.Code, bool, *Error) {
			return setters(gen, ctx, assignTo, sourceID, source, target, errPath, additionalFieldSources, definedFields, missing)
		})
		if ctx.Conf.SettersPreferred {
			fieldGens[0], fieldGens[1] = fieldGens[1], fieldGens[0]
		}
	}

	usedSourceID := false
	for _, fieldGen := range fieldGens {
		genStmt, genUsedSourceID, err := fieldGen()
		if err != nil {
			return nil, err
		}
		stmt = append(stmt, genStmt...)
		usedSourceID = usedSourceID || genUsedSourceID
	}

	for _, err := range missing {
		if err != nil {
			return nil, err
		}
	}

	if !usedSourceID {
		stmt = append(stmt, jen.Id("_").Op("=").Add(sourceID.Code.Clone()))
	}

	for name := range definedFields {
		return nil, NewError(fmt.Sprintf("Field %q does not exist.\nRemove or adjust field settings referencing this field.", name)).Lift(&Path{
			Prefix:     ".",
			TargetID:   name,
			TargetType: "???",
		})
	}

	return stmt, nil
}

func fields(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath, additionalFieldSources []xtype.FieldSources, definedFields map[string]struct{}, missing map[string]*Error) ([]jen.Code, bool, *Error) {
	stmt := []jen.Code{}

	usedSourceID := false
	for i := 0; i < target.StructType.NumFields(); i++ {
		targetField := target.StructType.Field(i)
		fieldStmt, fieldUsedSourceID, err := fieldGen(gen, ctx, assignTo, sourceID, source, target, errPath, additionalFieldSources, definedFields, missing, targetField)
		if err != nil {
			return nil, false, err
		}
		usedSourceID = usedSourceID || fieldUsedSourceID
		stmt = append(stmt, fieldStmt...)
	}
	return stmt, usedSourceID, nil
}

func setters(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath, additionalFieldSources []xtype.FieldSources, definedFields map[string]struct{}, missing map[string]*Error) ([]jen.Code, bool, *Error) {
	targetType := target.T

	stmt := []jen.Code{}
	usedSourceID := false

	if targetType, ok := targetType.(*types.Named); ok {
		methodSet := types.NewMethodSet(types.NewPointer(targetType))
		for i := 0; i < methodSet.Len(); i++ {
			method := methodSet.At(i).Obj().(*types.Func)
			targetField := types.NewVar(token.NoPos, nil, method.Name(), method.Type())
			fieldStmt, fieldUsedSourceID, err := fieldGen(gen, ctx, assignTo, sourceID, source, target, errPath, additionalFieldSources, definedFields, missing, targetField)
			if err != nil {
				return nil, false, err
			}
			usedSourceID = usedSourceID || fieldUsedSourceID
			stmt = append(stmt, fieldStmt...)
		}
	}
	return stmt, usedSourceID, nil
}

func fieldGen(gen Generator, ctx *MethodContext, assignTo *AssignTo, sourceID *xtype.JenID, source, target *xtype.Type, errPath ErrorPath, additionalFieldSources []xtype.FieldSources, definedFields map[string]struct{}, missing map[string]*Error, targetField *types.Var) ([]jen.Code, bool, *Error) {
	delete(definedFields, targetField.Name())

	fieldMapping := ctx.Field(target, targetField.Name())

	if fieldMapping.Ignore {
		return nil, false, nil
	}
	if !targetField.Exported() && ctx.Conf.IgnoreUnexported {
		return nil, false, nil
	}

	if !xtype.Accessible(targetField, ctx.OutputPackagePath) {
		cause := unexportedStructError(targetField.Name(), source.String, target.String)
		return nil, false, NewError(cause).Lift(&Path{
			Prefix:     ".",
			SourceID:   "???",
			TargetID:   targetField.Name(),
			TargetType: targetField.Type().String(),
		})
	}

	implicitName := targetField.Name()
	tempStmt := []jen.Code{}
	targetFieldType := xtype.TypeOf(targetField.Type())
	targetFieldPath := errPath.Field(targetField.Name())
	assignVar := assignTo.Stmt.Clone().Dot(targetField.Name())
	setterStmt := []jen.Code{}

	setterSignature, ok := targetField.Type().(*types.Signature)
	if ok && setterSignature.Recv() != nil {
		matches := ctx.Conf.SettersRegex.FindStringSubmatch(targetField.Name())
		if len(matches) == 2 && matches[0] == targetField.Name() {
			implicitName = matches[1]
		} else if fieldMapping.Source == "" && fieldMapping.Function == nil {
			return nil, false, nil
		}

		if setterSignature.Params().Len() != 1 {
			return nil, false, NewError(fmt.Sprintf("Setter method %s must have exactly one parameter.", targetField.Name())).Lift(&Path{
				Prefix:     ".",
				SourceID:   "???",
				TargetID:   targetField.Name(),
				TargetType: "???",
			})
		}

		setParam := setterSignature.Params().At(0)
		targetField := types.NewVar(token.NoPos, nil, targetField.Name(), setParam.Type())
		targetFieldType = xtype.TypeOf(targetField.Type())
		var err *Error
		tempStmt, assignVar, err = buildTargetVar(gen, ctx, sourceID, source, targetFieldType, targetFieldPath)
		if err != nil {
			return nil, false, err
		}
		setterStmt = []jen.Code{assignTo.Stmt.Clone().Dot(targetField.Name()).Call(assignVar.Clone())}
	}

	if err, ok := missing[implicitName]; ok && err == nil {
		return nil, false, nil
	}
	missing[implicitName] = nil

	stmt := []jen.Code{}
	usedSourceID := false
	if fieldMapping.Function == nil {
		nextID, nextSource, mapStmt, lift, miss, err := mapField(gen, ctx, implicitName, targetField, sourceID, source, target, additionalFieldSources, targetFieldPath)
		if miss {
			if ctx.Conf.IgnoreMissing {
				return nil, false, nil
			}
			storedErr := missing[implicitName]
			if storedErr != nil {
				return nil, false, storedErr
			} else {
				missing[implicitName] = err
				return nil, false, nil
			}
		}
		if err != nil {
			return nil, false, err
		}
		usedSourceID = true
		stmt = append(stmt, mapStmt...)

		fieldStmt := []jen.Code{}
		fieldStmt = append(fieldStmt, tempStmt...)

		assignStmt, err := gen.Assign(ctx, AssignOf(assignVar), nextID, nextSource, targetFieldType, targetFieldPath)
		if err != nil {
			return nil, false, err.Lift(lift...)
		}
		fieldStmt = append(fieldStmt, assignStmt...)

		fieldStmt = append(fieldStmt, setterStmt...)

		if shouldCheckAgainstZero(ctx, nextSource, targetFieldType, assignTo.Update, false) {
			stmt = append(stmt, jen.If(nextID.Code.Clone().Op("!=").Add(xtype.ZeroValue(nextSource.T))).Block(fieldStmt...))
		} else {
			stmt = append(stmt, fieldStmt...)
		}
	} else {
		def := fieldMapping.Function

		sourceLift := []*Path{}
		var functionCallSourceID *xtype.JenID
		var functionCallSourceType *xtype.Type
		if def.Source != nil {
			usedSourceID = true
			nextID, nextSource, mapStmt, mapLift, _, err := mapField(gen, ctx, implicitName, targetField, sourceID, source, target, additionalFieldSources, targetFieldPath)
			if err != nil {
				return nil, false, err
			}
			sourceLift = mapLift
			stmt = append(stmt, mapStmt...)

			if fieldMapping.Source == "." && sourceID.ParentPointer != nil &&
				def.Source.AssignableTo(source.AsPointer()) {
				functionCallSourceID = sourceID.ParentPointer
				functionCallSourceType = source.AsPointer()
			} else {
				functionCallSourceID = nextID
				functionCallSourceType = nextSource
			}
		} else {
			sourceLift = append(sourceLift, &Path{
				Prefix:     ".",
				TargetID:   targetField.Name(),
				TargetType: targetFieldType.String,
			})
		}

		fieldStmt := []jen.Code{}
		fieldStmt = append(fieldStmt, tempStmt...)

		callStmt, callReturnID, err := gen.CallMethod(ctx, fieldMapping.Function, functionCallSourceID, functionCallSourceType, targetFieldType, targetFieldPath)
		if err != nil {
			return nil, false, err.Lift(sourceLift...)
		}
		fieldStmt = append(fieldStmt, callStmt...)

		fieldStmt = append(fieldStmt, assignVar.Clone().Op("=").Add(callReturnID.Code))

		fieldStmt = append(fieldStmt, setterStmt...)

		if shouldCheckAgainstZero(ctx, functionCallSourceType, targetFieldType, assignTo.Update, true) {
			stmt = append(stmt, jen.If(functionCallSourceID.Code.Clone().Op("!=").Add(xtype.ZeroValue(functionCallSourceType.T))).Block(fieldStmt...))
		} else {
			stmt = append(stmt, fieldStmt...)
		}
	}
	return stmt, usedSourceID, nil
}

func shouldCheckAgainstZero(ctx *MethodContext, s, t *xtype.Type, isUpdate, call bool) bool {
	switch {
	case !ctx.Conf.UpdateTarget && !isUpdate:
		return false
	case s.Struct && ctx.Conf.IgnoreStructZeroValueField:
		return true
	case s.Basic && ctx.Conf.IgnoreBasicZeroValueField:
		return true
	case ctx.Conf.IgnoreNillableZeroValueField:
		if s.Chan || s.Map || s.Func || s.Signature || s.Interface {
			return true
		}
		if call || (ctx.Conf.SkipCopySameType && types.Identical(s.T, t.T)) {
			return (s.List && !s.ListFixed) || s.Pointer
		}
		return false
	default:
		return false
	}
}

func mapField(
	gen Generator,
	ctx *MethodContext,
	implicitName string,
	targetField *types.Var,
	sourceID *xtype.JenID,
	source, target *xtype.Type,
	additionalFieldSources []xtype.FieldSources,
	errPath ErrorPath,
) (*xtype.JenID, *xtype.Type, []jen.Code, []*Path, bool, *Error) {
	lift := []*Path{}
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
		return sourceID, source, nil, lift, false, nil
	}

	var path []string
	if pathString == "" {
		findField := []string{implicitName}
		if ctx.Conf.GettersEnabled {
			buf := &bytes.Buffer{}
			err := ctx.Conf.GettersTemplate.Execute(buf, implicitName)
			if err != nil {
				return nil, nil, nil, nil, false, NewError("Cannot execute getter template").Lift(&Path{
					Prefix:     ".",
					SourceID:   "???",
					TargetID:   targetField.Name(),
					TargetType: targetField.Type().String(),
				})
			}
			findField = append(findField, buf.String())
			if ctx.Conf.GettersPreferred {
				findField[0], findField[1] = findField[1], findField[0]
			}
		}

		var sourceMatch *xtype.StructField
		var err error
		for _, field := range findField {
			var ferr error
			sourceMatch, ferr = xtype.FindField(field, ctx.Conf.MatchIgnoreCase, source, additionalFieldSources)
			if ferr != nil {
				if err == nil {
					err = ferr
				}
			} else {
				err = nil
				break
			}
		}
		if err != nil {
			cause := fmt.Sprintf("Cannot match the target field with the source entry: %s.", err.Error())
			_, miss := err.(*xtype.NoMatchError)
			return nil, nil, nil, nil, miss, NewError(cause).Lift(&Path{
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

	nextIDCode := sourceID.Code
	nextSource := source

	for i := 0; i < len(path); i++ {
		if nextSource.Pointer {
			addCondition := nextIDCode.Clone().Op("!=").Nil()
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
			nextIDCode = nextIDCode.Clone().Dot(sourceMatch.Name)
			liftPath := &Path{
				Prefix:     ".",
				SourceID:   sourceMatch.Name,
				SourceType: nextSource.String,
			}

			if i == len(path)-1 {
				liftPath.TargetID = targetField.Name()
				liftPath.TargetType = targetField.Type().String()
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

	returnID := xtype.VariableID(nextIDCode)
	innerStmt := []jen.Code{}
	if nextSource.Func {
		def, err := method.Parse(nextSource.FuncType, &method.ParseOpts{
			Converter:         nil,
			OutputPackagePath: ctx.OutputPackagePath,
			ErrorPrefix:       "Error parsing struct method",
			Params:            method.ParamsNone,
			ContextMatch:      config.StructMethodContextRegex,
			CustomCall:        nextIDCode,
		}, method.EmptyLocalOpts)
		if err != nil {
			return nil, nil, nil, nil, false, NewError(err.Error()).Lift(lift...)
		}

		methodCallInner, callID, callErr := gen.CallMethod(ctx, def, nil, nil, def.Target, errPath)
		if callErr != nil {
			return nil, nil, nil, nil, false, callErr.Lift(lift...)
		}
		innerStmt = methodCallInner
		nextSource = def.Target
		returnID = callID
		lift = append(lift, &Path{
			Prefix:     "(",
			SourceID:   ")",
			SourceType: def.Target.String,
		})
	}

	if condition != nil && !nextSource.Pointer {
		lift[len(lift)-1].SourceType = fmt.Sprintf("*%s (It is a pointer because the nested property in the goverter:map was a pointer)",
			lift[len(lift)-1].SourceType)
	}

	stmt := []jen.Code{}
	if condition != nil {
		pointerNext := nextSource
		if !nextSource.Pointer {
			pointerNext = nextSource.AsPointer()
		}
		tempName := ctx.Name(pointerNext.ID())
		stmt = append(stmt, jen.Var().Id(tempName).Add(pointerNext.TypeAsJen()))

		if nextSource.Pointer {
			innerStmt = append(innerStmt, jen.Id(tempName).Op("=").Add(returnID.Code))
		} else {
			pstmt, pointerID := returnID.Pointer(nextSource, ctx.Name)
			innerStmt = append(innerStmt, pstmt...)
			innerStmt = append(innerStmt, jen.Id(tempName).Op("=").Add(pointerID.Code))
		}

		stmt = append(stmt, jen.If(condition).Block(innerStmt...))
		nextSource = pointerNext
		returnID = xtype.VariableID(jen.Id(tempName))
	} else {
		stmt = append(stmt, innerStmt...)
	}

	return returnID, nextSource, stmt, lift, false, nil
}

func parseAutoMap(ctx *MethodContext, source *xtype.Type) ([]xtype.FieldSources, *Error) {
	fieldSources := []xtype.FieldSources{}
	for _, field := range ctx.Conf.AutoMap {
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
				SourceType: field.Type.String,
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

See https://goverter.jmattheis.de/guide/unexported-field`, targetField)
}

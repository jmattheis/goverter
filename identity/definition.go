package identity

import (
	"strings"

	"github.com/jmattheis/goverter/xtype"
)

type Definition struct {
	OriginID string
	ID       string
	Package  string
	Name     string

	Type *xtype.Type
}

func (def *Definition) ArgDebug(indent string) string {
	var lines []string
	// for _, arg := range def.RawArgs {
	// 	argUse := arg.Use
	// 	if arg.Use == ArgUseMultiSource {
	// 		argUse = ArgUseSource
	// 	} else if arg.Use == ArgUseInterface {
	// 		argUse = ArgUseContext
	// 	}
	// 	lines = append(lines, fmt.Sprintf("[%s] %s", argUse, arg.Type.String))
	// }
	//
	// if def.Target != nil && !def.UpdateTarget {
	// 	lines = append(lines, fmt.Sprintf("[target] %s", def.Target.String))
	// }
	//
	// if len(lines) == 0 {
	// 	return ""
	// }

	return "\n" + indent + strings.Join(lines, "\n"+indent)
}

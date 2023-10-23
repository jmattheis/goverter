package generator

import (
	"fmt"

	"github.com/jmattheis/goverter/xtype"
)

func validateMethods(lookup map[xtype.Signature]*generatedMethod) error {
	for _, genMethod := range lookup {
		if genMethod.Explicit && len(genMethod.Fields) > 0 {
			isTargetStructPointer := genMethod.Target.Pointer && genMethod.Parameters.Target.PointerInner.Struct
			if !genMethod.Target.Struct && !isTargetStructPointer {
				return fmt.Errorf("Invalid struct field mapping on method:\n    %s\n\nField mappings like goverter:map or goverter:ignore may only be set on struct or struct pointers.\nSee https://goverter.jmattheis.de/#/config/nested", genMethod.ID)
			}
			if overlapping, ok := findOverlappingExplicitStructMethod(lookup, genMethod); ok {
				return fmt.Errorf("Overlapping field mapping found.\n\nPlease move the field related goverter:* comments from this method:\n    %s\n\nto this method:\n    %s\n\nGoverter will use %s inside the implementation of %s, thus, field related settings would be ignored.", genMethod.ID, overlapping.ID, overlapping.Name, genMethod.Name)
			}
		}
	}
	return nil
}

func findOverlappingExplicitStructMethod(lookup map[xtype.Signature]*generatedMethod, def *generatedMethod) (*generatedMethod, bool) {
	source := def.Source
	target := def.Target

	switch {
	case source.Struct && target.Pointer && target.PointerInner.Struct:
		genMethod, ok := lookup[xtype.SignatureOf(source, target.PointerInner)]
		if ok && genMethod.Explicit {
			return genMethod, true
		}
	case source.Pointer && source.PointerInner.Struct && target.Pointer && target.PointerInner.Struct:
		genMethod, ok := lookup[xtype.SignatureOf(source.PointerInner, target)]
		if ok && genMethod.Explicit {
			return genMethod, true
		}
		genMethod, ok = lookup[xtype.SignatureOf(source.PointerInner, target.PointerInner)]
		if ok && genMethod.Explicit {
			return genMethod, true
		}
	}
	return nil, false
}

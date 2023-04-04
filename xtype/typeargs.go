//go:build go1.18
// +build go1.18

package xtype

import "go/types"

func getTypeArgs(named *types.Named) []types.Type {
	result := []types.Type{}

	args := named.TypeArgs()
	for i := 0; i < args.Len(); i++ {
		result = append(result, args.At(i))
	}
	return result
}

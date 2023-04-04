//go:build !go1.18
// +build !go1.18

package xtype

import "go/types"

func getTypeArgs(named *types.Named) []types.Type {
	return nil
}

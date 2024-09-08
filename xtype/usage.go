package xtype

import (
	"sort"
)

func UsageFromMap[V any](value map[string]V) UsageChecker {
	m := map[string]struct{}{}

	for key := range value {
		m[key] = struct{}{}
	}

	return UsageChecker(m)
}

type UsageChecker map[string]struct{}

func (u UsageChecker) Used(key string) {
	delete(u, key)
}

func (u UsageChecker) Unused() []string {
	var keys []string
	for key := range u {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

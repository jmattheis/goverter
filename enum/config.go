package enum

import (
	"go/types"
	"regexp"
)

type Config struct {
	Unknown  string
	Enabled  bool
	Excludes IDPatterns
	// Package loads the package with the given path in full. Enum detection uses it for a
	// package known only through another package's export data, whose scope lacks the
	// constants that make a type an enum.
	Package func(path string) *types.Package
}

type IDPattern struct {
	Path *regexp.Regexp
	Name *regexp.Regexp
}

type IDPatterns []IDPattern

func (ids IDPatterns) Matches(path, name string) bool {
	for _, id := range ids {
		if id.Path.MatchString(path) && id.Name.MatchString(name) {
			return true
		}
	}
	return false
}

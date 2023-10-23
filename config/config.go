package config

import (
	"fmt"
	"go/types"
	"sort"

	"github.com/jmattheis/goverter/pkgload"
)

type RawLines struct {
	Location string
	Lines    []string
}

type RawConverter struct {
	Scope         *types.Scope
	InterfaceName string
	Converter     RawLines
	Methods       map[string]RawLines
	FileSource    string
}

type Raw struct {
	Loader     *pkgload.PackageLoader
	Converters []RawConverter
	Global     RawLines
}

func Parse(raw *Raw) ([]*Converter, error) {
	global, err := parseGlobal(raw.Loader, raw.Global)
	if err != nil {
		return nil, err
	}

	converters := []*Converter{}
	for _, rawConverter := range raw.Converters {
		converter, err := parseConverter(raw.Loader, &rawConverter, *global)
		if err != nil {
			return nil, err
		}
		converters = append(converters, converter)
	}

	sort.Slice(converters, func(i, j int) bool {
		return converters[i].Name < converters[j].Name
	})

	return converters, nil
}

func formatLineError(lines RawLines, t, value string, err error) error {
	cmd, _ := parseCommand(value)
	msg := `error parsing 'goverter:%s' at
    %s
    %s

%s`
	return fmt.Errorf(msg, cmd, lines.Location, t, err)
}

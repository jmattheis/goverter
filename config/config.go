package config

import (
	"fmt"
	"sort"

	"github.com/jmattheis/goverter/config/parse"
	"github.com/jmattheis/goverter/enum"
	"github.com/jmattheis/goverter/pkgload"
)

type RawLines struct {
	Location string
	Lines    []string
}

type RawConverter struct {
	PackagePath   string
	PackageName   string
	InterfaceName string
	Converter     RawLines
	Methods       map[string]RawLines
	FileName      string
}

type Raw struct {
	Converters []RawConverter
	Global     RawLines

	WorkDir              string
	BuildTags            string
	OuputBuildConstraint string

	EnumTransformers map[string]enum.Transformer
}

type context struct {
	Loader           *pkgload.PackageLoader
	WorkDir          string
	EnumTransformers map[string]enum.Transformer
}

func Parse(raw *Raw) ([]*Converter, error) {
	loader, err := pkgload.New(raw.WorkDir, raw.BuildTags, getPackages(raw))
	if err != nil {
		return nil, err
	}

	ctx := &context{Loader: loader, EnumTransformers: raw.EnumTransformers, WorkDir: raw.WorkDir}

	converters := []*Converter{}
	for _, rawConverter := range raw.Converters {
		converter, err := parseConverter(ctx, &rawConverter, raw.Global)
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
	cmd, _ := parse.Command(value)
	msg := `error parsing 'goverter:%s' at
    %s
    %s

%s`
	return fmt.Errorf(msg, cmd, lines.Location, t, err)
}

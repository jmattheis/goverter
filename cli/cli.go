package cli

import (
	"flag"
	"fmt"

	"github.com/jmattheis/goverter"
	"github.com/jmattheis/goverter/config"
)

type Strings []string

func (s Strings) String() string {
	return fmt.Sprint([]string(s))
}

func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func Parse(args []string) (*goverter.GenerateConfig, error) {
	switch len(args) {
	case 0:
		return nil, usage("unknown")
	case 1:
		return nil, usageErr("missing command: ", args[0])
	default:
		if args[1] != "gen" {
			return nil, usageErr("unknown command: "+args[1], args[0])
		}
	}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.Usage = func() {
		// fmt.Println(usage(args[0]))
	}

	var global Strings
	fs.Var(&global, "global", "add ")
	fs.Var(&global, "g", "add ")

	if err := fs.Parse(args[2:]); err != nil {
		return nil, usageErr(err.Error(), args[0])
	}
	patterns := fs.Args()

	if len(patterns) == 0 {
		return nil, usageErr("missing PATTERN", args[0])
	}

	c := goverter.GenerateConfig{
		PackagePatterns: patterns,
		Global: config.RawLines{
			Lines:    global,
			Location: "command line (-g, -global)",
		},
	}

	return &c, nil
}

func usageErr(err, cmd string) error {
	return fmt.Errorf("Error: %s\n%s", err, usage(cmd))
}

func usage(cmd string) error {
	return fmt.Errorf(`Usage:
  %s gen [OPTIONS] PACKAGE...

PACKAGE(s):
  Define the import paths goverter will use to search for converter interfaces.
  You can define multiple packages and use the special ... golang pattern to
  select multiple packages. See $ go help packages

OPTIONS:
  -g [value], -global [value]:
      apply settings to all defined converters. For a list of available
      settings see: https://goverter.jmattheis.de/#/config/

  -h, --help:
      display this help page

Examples:
  %s gen ./example/simple ./example/complex
  %s gen ./example/...
  %s gen github.com/jmattheis/goverter/example/simple
  %s gen -g ignoreMissing -g 'output ./generated/generated.go'  ./simple

Output:
  The output setting is relative to the conversion interface. If you want it
  relative to the current working directory you can use the magic @cwd path.

  %s gen -g 'output @cwd/generated'  ./simple

Documentation:
  Full documentation is available here: https://goverter.jmattheis.de`, cmd, cmd, cmd, cmd, cmd, cmd)
}

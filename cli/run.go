package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jmattheis/goverter"
	"github.com/jmattheis/goverter/enum"
)

type RunOpts struct {
	EnumTransformers map[string]enum.Transformer
}

// Run runs the goverter cli with the given args and customizations.
func Run(args []string, opts RunOpts) {
	cmd, err := Parse(args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch cmd := cmd.(type) {
	case *Help:
		_, _ = fmt.Fprintln(os.Stdout, cmd.Usage)
		os.Exit(0)
	case *Generate:
		if opts.EnumTransformers != nil {
			for key, value := range opts.EnumTransformers {
				cmd.Config.EnumTransformers[key] = value
			}
		}

		if err = goverter.GenerateConverters(cmd.Config); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *Version:
		b, ok := debug.ReadBuildInfo()
		if ok {
			fmt.Println(b)
		}
	default:
		panic("unknown command")
	}
}

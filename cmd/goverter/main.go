package main

import (
	"fmt"
	"os"

	goverter "github.com/jmattheis/goverter"
	"github.com/jmattheis/goverter/cli"
)

func main() {
	cfg, err := cli.Parse(os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = goverter.GenerateConverters(cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	goverter "github.com/sergeydobrodey/goverter"
)

func main() {
	packageName := flag.String("packageName", "generated", "")
	output := flag.String("output", "./generated/generated.go", "")
	extends := flag.String("extends", "", "comma separated list of local or package extends")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "expected one argument")
		return
	}
	pattern := args[0]
	var extendMethods []string
	if *extends != "" {
		extendMethods = strings.Split(*extends, ",")
	}

	err := goverter.GenerateConverterFile(*output, goverter.GenerateConfig{
		PackageName:   *packageName,
		ScanDir:       pattern,
		ExtendMethods: extendMethods,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}

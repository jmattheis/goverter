package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	goverter "github.com/jmattheis/goverter"
)

func main() {
	packageName := flag.String("packageName", "generated", "")
	output := flag.String("output", "./generated/generated.go", "")
	extends := flag.String("extends", "", "comma separated list of local or package extends")
	packagePath := flag.String("packagePath", "", "optional full package path for the generated code")
	wrapErrors := flag.Bool("wrapErrors", false,
		"if set, wrap conversion errors with extra details, such as struct field names")

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
		PackagePath:   *packagePath,
		WrapErrors:    *wrapErrors,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

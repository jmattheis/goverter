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
	workingDir := flag.String("workingDir", "", "optional working directory, default is the current directory")
	packagePath := flag.String("packagePath", "", "optional full package path for the generated code")

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

	if *workingDir != "" {
		if err := os.Chdir(*workingDir); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "could not change directory:", err)
			return
		}
	}

	err := goverter.GenerateConverterFile(*output, goverter.GenerateConfig{
		PackageName:   *packageName,
		ScanDir:       pattern,
		ExtendMethods: extendMethods,
		WorkingDir:    *workingDir,
		PackagePath:   *packagePath,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}

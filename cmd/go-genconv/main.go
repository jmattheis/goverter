package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmattheis/go-genconv"
)

func main() {
	packageName := *flag.String("packageName", "generated", "")
	output := *flag.String("output", "./generated/generated.go", "")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "expected one argument")
		return
	}
	pattern := args[0]

	err := genconv.GenerateConverterFile(output, genconv.GenerateConfig{
		PackageName: packageName,
		ScanDir:     pattern,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}

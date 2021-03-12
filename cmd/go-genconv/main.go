package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmattheis/go-genconv"
)

func main() {
	packageName := *flag.String("packageName", "generated", "")
	output := *flag.String("output", "./generated/generated.go", "")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("expected one argument")
		return
	}
	pattern := args[0]

	file, err := genconv.Generate(genconv.GenerateConfig{
		PackageName: packageName,
		ScanDir:     pattern,
	})
	os.MkdirAll(path.Dir(output), 0755)

	err = ioutil.WriteFile(output, file, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

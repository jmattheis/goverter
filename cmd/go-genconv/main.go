package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmattheis/go-genconv/internal/comments"
	"github.com/jmattheis/go-genconv/internal/generator"
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

	mapping, err := comments.ParseDocs(pattern)
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := generator.Generate(pattern, mapping, generator.Config{Name: packageName})
	if err != nil {
		fmt.Println(err)
		return
	}

	os.MkdirAll(path.Dir(output), 0755)
	buf := &bytes.Buffer{}
	err = file.Render(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(output, buf.Bytes(), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

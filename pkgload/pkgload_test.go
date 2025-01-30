package pkgload

import (
	"fmt"
)

func ExamplePackageLoader_LoadPkgPathFromDir() {
	loader, err := New(".", "", []string{})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	name, path, err := loader.LoadPkgPathFromDir(".")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("name: %q\n", name)
	fmt.Printf("path: %q\n", path)
	// Output:
	// name: "pkgload"
	// path: "github.com/jmattheis/goverter/pkgload"
}

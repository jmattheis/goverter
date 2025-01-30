package pkgload

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
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

func TestLoadPkgPathFromDir(t *testing.T) {
	absDir, err := filepath.Abs(".")
	require.NoError(t, err, "abs")
	t.Logf("abs dir path: %s", absDir)

	t.Run("uses cache", func(t *testing.T) {
		loader, err := New("", "", []string{})
		require.NoError(t, err, "new")

		want := &packages.Package{
			Name:    "my-pkg",
			PkgPath: "github.com/jmattheis/goverter/my-pkg",
		}
		loader.lookupAbsDir[absDir] = want

		name, pkgPath, err := loader.LoadPkgPathFromDir(absDir)
		require.NoError(t, err, "load")

		assert.Equal(t, "my-pkg", name, "package name")
		assert.Equal(t, "github.com/jmattheis/goverter/my-pkg", pkgPath, "package path")
		assert.Same(t, want, loader.lookupAbsDir[absDir], "cached package pointers")
	})

	t.Run("loads new pkg", func(t *testing.T) {
		loader, err := New("", "", []string{})
		require.NoError(t, err, "new")

		delete(loader.lookupAbsDir, absDir)

		name, pkgPath, err := loader.LoadPkgPathFromDir(absDir)
		require.NoError(t, err, "load")

		assert.Equal(t, "pkgload", name, "package name")
		assert.Equal(t, "github.com/jmattheis/goverter/pkgload", pkgPath, "package path")
		assert.NotNil(t, loader.lookupAbsDir[absDir], "cached package pointers")
	})
}

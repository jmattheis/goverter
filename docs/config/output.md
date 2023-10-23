## output:file

`output:file FILE` can be defined as [CLI argument](config/define.md#cli) or
[converter comment](config/define.md#converter). Default is
`./generated/generated.go`.

`output:file` sets the generated output file of a converter interface.
The file location is relative to the file with the converter interface. E.g. if
the interface is located at `/home/jm/src/pkg/converter/interface.go` and you
define `output:file ../generated/output.go` then the file would be created at
`/home/jm/src/pkg/generated/interface.go`

You can use the magic `@cwd/` prefix to reference the current working directory
goverter is executed in. E.g. you execute goverter in `/home/jm/src/pkg/,` the
interface is located at `/home/jm/src/pkg/converter/interface.go` and
`output:file @cwd/output/generated.go` is defined then the file would be
created at `/home/jm/src/pkg/output/generated.go`.

Different converters may have the same `output:file` if the `output:package` is
the same. See this more complex example:

<!-- tabs:start -->

#### **root.go**

```go
package example

// goverter:converter
// goverter:output:file ./a/generated.go
// goverter:output:package goverter/example/a
type RootA interface{
    Convert([]bool) []bool
}

// goverter:converter
// goverter:output:file ./b/generated.go
// goverter:output:package goverter/example/b
type RootB interface {
    Convert([]string) []string
}
```

#### **c/input.go**

```go
package c

// goverter:converter
// goverter:output:file ../a/generated.go
// goverter:output:package goverter/example/a
type CIntoA interface {
    Convert([]int) []int
}
```

#### **a/generated.go**

```go
package a

type CIntoAImpl struct{}

func (c *CIntoAImpl) Convert(source []int) []int {
	var intList []int
	if source != nil {
		intList = make([]int, len(source))
		for i := 0; i < len(source); i++ {
			intList[i] = source[i]
		}
	}
	return intList
}

type RootAImpl struct{}

func (c *RootAImpl) Convert(source []bool) []bool {
	var boolList []bool
	if source != nil {
		boolList = make([]bool, len(source))
		for i := 0; i < len(source); i++ {
			boolList[i] = source[i]
		}
	}
	return boolList
}
```


#### **b/generated.go**

```go
package b

type RootBImpl struct{}

func (c *RootBImpl) Convert(source []string) []string {
	var stringList []string
	if source != nil {
		stringList = make([]string, len(source))
		for i := 0; i < len(source); i++ {
			stringList[i] = source[i]
		}
	}
	return stringList
}
```

<!-- tabs:end -->

## output:package

`output:package [PACKAGE][:NAME]` can be defined as
[CLI argument](config/define.md#cli) or
[converter comment](config/define.md#converter). Default is `:generated`.

### output:package PACKAGE

This is the recommended way to define this setting. If you define the full
package path, goverter is able to prevent edge-cases in the converter
generation. The package name used in the generated `.go` file will be inferred
from the normalized full package path. E.g.

```go
// goverter:converter
// goverter:output:package example.org/simple/mypackage
type Converter interface {
    Convert([]string) []string
}
```

will create  **generated/generated.go** starting with

```go
package mypackage
// ...
```

If there are reserved characters in the package name will be normalized. E.g.

```go
// goverter:converter
// goverter:output:package example.org/simple/my-cool_package
type Converter interface {
    Convert([]string) []string
}
```

will create  **generated/generated.go** starting with

```go
package mycoolpackage
// ...
```


### output:package PACKAGE:NAME

If you want to overwrite the inferred package name, you can do so with
`output:package PACKAGE:NAME`. E.g.

```go
// goverter:converter
// goverter:output:package example.org/simple/my-cool_package:overriddenname
type Converter interface {
    Convert([]string) []string
}
```

will create  **generated/generated.go** starting with

```go
package overriddenname
// ...
```


### output:package NAME

If you aren't able to define the full package path, then you can only define
the package name with `output:package :NAME`. Goverter may produce uncompilable
code, when you hit an edge-case that requires a full package path. If the
generated code compiles, then there should be no problems using the setting
like this.

```go
// goverter:converter
// goverter:output:package :mypackage
type Converter interface {
    Convert([]string) []string
}
```

will create  **generated/generated.go** starting with

```go
package mypackage
// ...
```

# Setting: output

## output:file

`output:file FILE` can be defined as [CLI argument](./define-settings.md#cli) or
[converter comment](./define-settings.md#converter). Default is
`./generated/generated.go`.

`output:file` sets the generated output file of a converter interface. The file
location is relative to the file with the converter interface. E.g. if the
interface is located at `/home/jm/src/pkg/converter/interface.go` and you define
`output:file ../generated/output.go` then the file would be created at
`/home/jm/src/pkg/generated/interface.go`

You can use the magic `@cwd/` prefix to reference the current working directory
goverter is executed in. E.g. you execute goverter in `/home/jm/src/pkg/,` the
interface is located at `/home/jm/src/pkg/converter/interface.go` and
`output:file @cwd/output/generated.go` is defined then the file would be created
at `/home/jm/src/pkg/output/generated.go`.

Different converters may have the same `output:file` if the `output:package` is
the same. See this more complex example:

::: code-group
<<< @../../example/output-multiple-files/root.go
<<< @../../example/output-multiple-files/c/input.go [c/input.go]
<<< @../../example/output-multiple-files/a/generated.go [c/generated.go]
<<< @../../example/output-multiple-files/b/generated.go [b/generated.go]
:::

## output:package

`output:package [PACKAGE][:NAME]` can be defined as
[CLI argument](./define-settings.md#cli) or [converter comment](./define-settings.md#converter).
Default is `:generated`.

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

will create **generated/generated.go** starting with

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

will create **generated/generated.go** starting with

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

will create **generated/generated.go** starting with

```go
package overriddenname
// ...
```

### output:package NAME

If you aren't able to define the full package path, then you can only define the
package name with `output:package :NAME`. Goverter may produce uncompilable
code, when you hit an edge-case that requires a full package path. If the
generated code compiles, then there should be no problems using the setting like
this.

```go
// goverter:converter
// goverter:output:package :mypackage
type Converter interface {
    Convert([]string) []string
}
```

will create **generated/generated.go** starting with

```go
package mypackage
// ...
```

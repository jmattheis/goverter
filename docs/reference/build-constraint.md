# Build Constraint/Tags

Introduced in v1.3.0.

The defaults are configurable via the [cli](./cli.md).

By default, goverter will consider the build tag `goverter` satisfied during
the scanning & generation of conversion implementations and all generated files
define the [Build Constraint](https://pkg.go.dev/cmd/go#hdr-Build_constraints):

```go
//go:build !goverter
```

This constraint isn't satisfied during the goverter generation, therefore this
file is excluded and ignored. This feature is to ignore and suppress compile
errors in unneeded files, so goverter is able to fix them with the newly
generated implementation.

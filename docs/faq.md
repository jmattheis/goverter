# Frequently Asked Questions

## Generate code into the same package

Internally Goverter cannot automatically infer the target package path. This
path is required to correctly import relative types. To fix this you have to
configure the full package path in
[`output:package`](reference/output.md#outputpackage).

E.g.
```go
// goverter:converter
// goverter:package github.com/jmattheis/goverter/example/sample
type Converter interface {}
```

Afterwards, goverter should correctly import types in the same package.

## import cycle not allowed

See [Generate code into the same package](#generate-code-into-the-same-package)

## Generate only shallow copy

See [`skipCopySameType`](reference/skipCopySameType.md).

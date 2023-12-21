# Configure Nested

If you want to [`map`](../reference/map.md) or
[`ignore`](../reference/ignore.md) a field of a nested type like a map or
slice, then you've to define another converter method.

Example, you've want to `map` the `NestedInput.LastName` to
`NestedOutput.Surname` for this method.

::: code-group
<<< @../../example/nested-struct/model.go
:::

You can't do it like this:

```go
// goverter:converter
type Converter interface {
    // goverter:map Nested.LastName Nested.Surname
    Convert(Input) Output
}
```

because Goverter dosen't support nested target fields. You have to create
another conversion method targeting the nested types like this:

::: code-group
<<< @../../example/nested-struct/input.go
<<< @../../example/nested-struct/generated/generated.go [generated/generated.go]
<<< @../../example/nested-struct/model.go
:::

## Slices and Maps

The rule above applies to the conversion of slices and maps too. Field settings
must be directly on the struct. This e.g. would be error:

```go
// goverter:converter
type Converter interface {
    // goverter:map LastName Surname
    Convert([]NestedInput) []NestedOutput
}
```

and should be written as:
```go
// goverter:converter
type Converter interface {
    Convert([]NestedInput) []NestedOutput
	// goverter:map LastName Surname
	ConvertNested(NestedInput) NestedOutput
}
```

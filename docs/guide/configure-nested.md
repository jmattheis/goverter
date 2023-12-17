# Configure Nested

If you want to [`map`](../reference/map.md) or [`ignore`](../reference/ignore.md) on a
nested type like a map or slice, then you've to define another converter method.

Example, you've want to `map` the `NestedInput.LastName` to
`NestedOutput.Surname` for this method.

```go
package example

// goverter:converter
type Converter interface {
    Convert([]Input) []Output
}

type Input struct {
    Name string
    Nested NestedInput
}
type NestedInput struct {
    LastName string
}

type Output struct {
    Name string
    Nested NestedOutput
}
type NestedOutput struct {
    Surname string
}
```

You can't do it like this:

```go
// goverter:converter
type Converter interface {
    // goverter:map LastName Surname
    Convert([]Input) []Output
}
```

because Goverter doesn't doesn't apply the mapping to all sub conversions, as
this could cause unexpected behavior. The correct way would be to define another
conversion method like this:

::: code-group
<<< @../../example/nested-struct/input.go
<<< @../../example/nested-struct/generated/generated.go [generated/generated.go]
:::

# How to handle unexported fields

Goverter will fail when it sees unexported fields and requires the user to
define what to do. Here are ways to fix the unexported fields error in
Goverter:

[[toc]]

## Ignore the unexported field

If the unexported fields can be ignored, then you can use
[`ignore`](../reference/ignore.md) to ignore these fields.

An example for are protobuf generated structs, for which the fields `state`,
`sizeCache` and `unknownFields` are only internally used and must not be set
when constructing a new instance. It's possible to globally ignore unexported
fields via [`ignoreUnexported`](../reference/ignoreUnexported.md) though it's
_discouraged_ to use it.

::: details Example (click me)
::: code-group
<<< @../../example/protobuf/input.go
<<< @../../example/protobuf/event.proto
<<< @../../example/protobuf/pb/event.pb.go [pb/event.pb.go]
<<< @../../example/protobuf/generated/generated.go [generated/generated.go]
:::

## Shallow copy

If the unexported field is inside an immutable struct, you can skip the deep
copying by defining an [`extend`](../reference/extend) or [`map [SOURCE_PATH] TARGET
| METHOD`](../reference/map.md#map-SOURCE-PATH-TARGET-METHOD) method that
returns the input.

An example for this use-case would be the `time.Time` struct from the Go
standard library. The type is immutable and therefore can be safely passed
through.

::: details Example (click me)
::: code-group
<<< @../../example/time/input.go
<<< @../../example/time/generated/generated.go [generated/generated.go]
:::

## Use a constructor

If the struct containing the unexported fields can be instantiated via a
constructor you can use [`extend`](../reference/extend.md) or [`map
[SOURCE_PATH] TARGET | METHOD`](../reference/map.md#map-SOURCE-PATH-TARGET-METHOD)


::: details Example (click me)
::: code-group
<<< @../../example/constructor/input.go
<<< @../../example/constructor/generated/generated.go [generated/generated.go]
:::

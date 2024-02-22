# How to convert enums

[[toc]]

## Definition

The Go language doesn't explicit define enums. In the context of Goverter an
enum type is defined as a named type with an [underlying
type](https://go.dev/ref/spec#Underlying_types) of (`float`, `string`, or
`integer`) having at least one constant defined in the same package.

These examples would qualify as enums:

::: code-group
<<< @../../example/enum/unknown/input/enum.go [iota.go]
<<< @../../example/enum/unknown/output/enum.go [string.go]
:::

## Conversion

When goverter sees types which both qualify as enum according to the definition
above then it tries to convert each key of the source enum to the target enum.
When not all keys of the source enum can be mapped, goverter will error.

## Unknown enum values

In Go it's possible to manually instantiate values of types qualifying as
[Enum](#definition). E.g.

```go
func GetColor() Color { return Color("invalid") }

type Color string
const Green Color = "green"
const Blue  Color = "blue"
```

it's not clear how to handle the `Color("invalid")`, therefore _you_ have to
explicitly configure how to convert an invalid enum value. You can do this by
configuring [`enum:unknown`](../reference/enum.md#enum-unknown-action) with one
of these actions.

`enum:unknown @error`: return an error when an invalid value is encountered.

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/error/input.go
<<< @../../example/enum/unknown/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/error/generated/generated.go [generated/generated.go]
:::

`enum:unknown @ignore`: ignore invalid values and return the zero value of the enum type.

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/ignore/input.go
<<< @../../example/enum/unknown/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/ignore/generated/generated.go [generated/generated.go]
:::

`enum:unknown @panic`: panic when an invalid value is encountered

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/panic/input.go
<<< @../../example/enum/unknown/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/panic/generated/generated.go [generated/generated.go]
:::

`enum:unknown KEY`: use an enum key when an invalid value in encountered

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/key/input.go
<<< @../../example/enum/unknown/key/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/key/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/key/generated/generated.go [generated/generated.go]
:::

## Mapping enum keys

If your source and target enum have differently named keys, you can use
[`enum:map`](../reference/enum.md#enum-map-source-target) to define the mapping.

::: details Example (click me)
::: code-group
<<< @../../example/enum/map/input.go
<<< @../../example/enum/map/input/enum.go [input/enum.go]
<<< @../../example/enum/map/output/enum.go [output/enum.go]
<<< @../../example/enum/map/generated/generated.go [generated/generated.go]
:::

You can also use any of the `@actions` from [`enum:unknown`](#unknown-enum-values). E.g.

::: details Example (click me)
::: code-group
<<< @../../example/enum/map-panic/input.go
<<< @../../example/enum/map-panic/input/enum.go [input/enum.go]
<<< @../../example/enum/map-panic/output/enum.go [output/enum.go]
<<< @../../example/enum/map-panic/generated/generated.go [generated/generated.go]
:::

If you have keys that differ in the same format you can use the transformer
[`enum:transform regex`](../reference/enum.md#enum-transform-regex-search-replace)
to do a regex search and replace.

::: details Example (click me)
::: code-group
<<< @../../example/enum/transform-regex/input.go
<<< @../../example/enum/transform-regex/generated/generated.go [generated/generated.go]
:::

If that's still not powerful enough for you, then you can define [a custom enum
transformer](../reference/enum.md#enum-transform-custom)

## Disable enum detection and conversion

To disable enum detection and conversion, add `enum no` to the converter.

::: details Example (click me)
::: code-group
<<< @../../example/enum/disable/input.go
<<< @../../example/enum/disable/generated/generated.go [generated/generated.go]
:::

## Handle false positives

You can exclude wrongly detected enums via
[`enum:exclude`](../reference/enum.md#enum-exclude).

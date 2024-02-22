# Setting: enum

### Definition

The Go language doesn't explicit define enums. In the context of Goverter an
enum type is defined as a named type with an [underlying
type](https://go.dev/ref/spec#Underlying_types) of (`float`, `string`, or
`integer`) having at least one constant defined in the same package.

## enum

`enum [yes|no]` can be defined as [CLI argument](./define-settings.md#cli) or
[converter comment](./define-settings.md#converter). `enum` is enabled per
default.

`enum` allows you to disable enum support in goverter.

::: details Example (click me)
::: code-group
<<< @../../example/enum/disable/input.go
<<< @../../example/enum/disable/generated/generated.go [generated/generated.go]
:::

## enum:unknown ACTION


`enum:unknown ACTION|KEY` can be defined as [CLI
argument](./define-settings.md#cli), [converter
comment](./define-settings.md#converter) or [method
comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

Define what happens on an invalid or unexpected enum value.

`enum:unknown @error` returns an error in the default case of the switch
statement.

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/error/input.go
<<< @../../example/enum/unknown/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/error/generated/generated.go [generated/generated.go]
:::

`enum:unknown @ignore` does nothing in the default case of the switch
statement.

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/ignore/input.go
<<< @../../example/enum/unknown/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/ignore/generated/generated.go [generated/generated.go]
:::

`enum:unknown @panic` panics in the default case of the switch statement.

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/panic/input.go
<<< @../../example/enum/unknown/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/panic/generated/generated.go [generated/generated.go]
:::

`enum:unknown KEY`: use an enum key in the default case of the switch statement

::: details Example (click me)
::: code-group
<<< @../../example/enum/unknown/key/input.go
<<< @../../example/enum/unknown/key/input/enum.go [input/enum.go]
<<< @../../example/enum/unknown/key/output/enum.go [output/enum.go]
<<< @../../example/enum/unknown/key/generated/generated.go [generated/generated.go]
:::

## enum:exclude

`enum:exclude [PACKAGE:]NAME` can be defined as [CLI
argument](./define-settings.md#cli) or [converter
comment](./define-settings.md#converter).

You can exclude falsely detected enums with exclude. This is useful when a type
[qualifies as enum](#definition) but isn't one. If `PACKAGE` is unset, goverter
will use the package of the converter interface.

Both `PACKAGE` and `NAME` can be regular expressions.

::: details Example (click me)
::: code-group
<<< @../../example/enum/exclude/input.go
<<< @../../example/enum/exclude/generated/generated.go [generated/generated.go]
:::

## enum:map SOURCE TARGET

`enum:map SOURCE TARGET` can be defined as [method
comment](./define-settings.md#method).

`enum:map` can be used to map an enum key of the source enum type to the target
type. In this example we have two spellings of Gray|Grey.

::: details Example (click me)
::: code-group
<<< @../../example/enum/map/input.go
<<< @../../example/enum/map/input/enum.go [input/enum.go]
<<< @../../example/enum/map/output/enum.go [output/enum.go]
<<< @../../example/enum/map/generated/generated.go [generated/generated.go]
:::

You can also use any of the `@actions` from
[`enum:unknown`](#enum-unknown-action). E.g.

::: details Example (click me)
::: code-group
<<< @../../example/enum/map-panic/input.go
<<< @../../example/enum/map-panic/input/enum.go [input/enum.go]
<<< @../../example/enum/map-panic/output/enum.go [output/enum.go]
<<< @../../example/enum/map-panic/generated/generated.go [generated/generated.go]
:::

## enum:transform ID [CONFIG]

`enum:transform ID [CONFIG]` can be defined as [method
comment](./define-settings.md#method).

`enum:transform` allows you to transform multiple enum keys to the target keys.
There is currently one builtin transformers, but you can define transformers
yourself.

### enum:transform regex SEARCH REPLACE

With `regex` you can do search and replace with regex of the enum keys.

::: details Example (click me)
::: code-group
<<< @../../example/enum/transform-regex/input.go
<<< @../../example/enum/transform-regex/generated/generated.go [generated/generated.go]
:::

### enum:transform CUSTOM

You can define a custom transformer by adding [goverter as dependency to your
project ](/guide/install.md#dependency) and then creating a new `main` package in your project
for your customized goverter. E.g. `./goverter/run.go`. Afterwards, you have to
pass your custom enum transformer to `cmd.Run` and then you can execute it via
`go run ./goverter` (use the package where you created you customized
goverter). Here is an example that adds a trim-prefix transformer.

::: details Example (click me)
::: code-group
<<< @../../example/enum/transform-custom/goverter/run.go [./goverter/run.go]
<<< @../../example/enum/transform-custom/input.go
<<< @../../example/enum/transform-custom/generated/generated.go [generated/generated.go]
:::

::: details Regex Implementation (click me)
::: code-group
<<< @../../enum/transformer_builtin.go
:::

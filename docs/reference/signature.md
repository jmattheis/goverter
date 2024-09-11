# Signature

Goverter supports a variety of different signatures.

[[toc]]

## Definition

In Go a function has **params** and **results**. E.g.

```go
type Converter interface {
    Convert(A, B) (C, D)
}
```

- `A` and `B` are **param**s
- `C` and `C` are **result**s or returns

### Named params/results

Go allows you to specify names to **params** and **results**. This document
references them as "named params" / "named results".

```go
type Converter interface {
    Convert(nameA A, nameB B) (nameC C, nameD D)
}
```

## Categories

Goverter groups **params** and **results** into these categories:

1. `source` is a **param** that will be converted to the `target` type. Any
   **param** that is not a `context` type is a `source` type.
1. `target` is the **first result** which `sources` are converted to.
1. `target_error` is the **second result** type which can be specified if the
   conversion fails. `target_error` must be of type
   [`error`](https://go.dev/tour/methods/19)
1. `context` are **[named](#named-paramsresults) params** where the name starts
   with `ctx` or `context`. The regex can be adjusted via
   [`arg:context:regex`](./arg.md#arg-context-regex). They are used in
   [custom functions](#custom-function) for manual conversion. `context` types
   aren't used for automatic conversion.

### Default context

If you are using [`goverter:converter`](./converter.md), then the conversion
interface seen as `context` regardless of the param name.

## Supported Signatures

### Conversion Method

Depending on the [`output:format`](./output.md#output-format), a conversion
method is defined on an interface, e.g. `ConvertA` or as variable, e.g.
`ConvertB`:

```go
// goverter:converter
type Converter interface {
    ConvertA(Input) Output
}
// goverter:variables
var (
    ConvertB func(Input) Output
)
```

Conversion methods support the following categories:

- `source`: required
- `context`: optional (multiple), a conversion method may have zero or more
  `context` types
- `target`: required
- `target_error`: optional

Here are signatures that would satisfy the given requirements:

```go
// goverter:converter
type Converter interface {
    ConvertOne(A) B
    // A=source; B=target

    ConvertTwo(source A) (B, error)
    // A=source; B=target; error=target_error

    ConvertThree(source A, ctx B) C
    // A=source; B=context; B=target

    ConvertFour(ctxOne A, source B, ctxSecond C) (D, error)
    // A=context; B=source; C=context; D=target; error=target_error
}
```

::: details Example (click to expand)
::: code-group
<<< @../../example/context/database/input.go
<<< @../../example/context/database/generated/generated.go [generated/generated.go]
:::

### Signature: optional source

The settings [`default`](./default.md) and
[`map CUSTOM`](./map.md#map-source-path-target-method) support the following
categories:

- `source`: optional
- `context`: optional (multiple), may have zero or more `context` types
- `target`: required
- `target_error`: optional

Here are signatures that would satisfy the given requirements:

```go
func ConvertOne() A {/**/}
// A=target

func ConvertTwo() (B, error) {/**/}
// A=target; error=target_error

func ConvertThree(input A) B {/**/}
// A=source; B=target

func ConvertFour(input A) (B, error) {/**/}
// A=source; B=target; error=target_error

func ConvertFour(ctxOne A, ctxTwo B) C {/**/}
// A=context; B=context; C=target

func ConvertFive(input A, ctxOne B, ctxTwo C) D {/**/}
// A=source; B=context; C=context; D=target

func ConvertSix(ctxOne A, source B ctxTwo C) (D, error) {/**/}
// A=context; B=source; C=context; D=target; error=target_error
```

::: details Example (click to expand)
::: code-group
<<< @../../example/context/database/input.go
<<< @../../example/context/database/generated/generated.go [generated/generated.go]
:::

### Signature: required source

[`extend`](./extend.md) supports the following categories:

- `source`: required
- `context`: optional (multiple), may have zero or more `context` types
- `target`: required
- `target_error`: optional

Here are signatures that would satisfy the given requirements:

```go
func ConvertOne(input A) B {/**/}
// A=source; B=target

func ConvertTwo(input A) (B, error) {/**/}
// A=source; B=target; error=target_error

func ConvertThree(input A, ctxOne B) C {/**/}
// A=source; B=context; C=target

func ConvertFour(input A, ctxOne B, ctxTwo C) D {/**/}
// A=source; B=context; C=context; D=target

func ConvertFive(ctxOne A, source B ctxTwo C) (D, error) {/**/}
// A=context; B=source; C=context; D=target; error=target_error
```

::: details Example (click to expand)
::: code-group
<<< @../../example/context/date-format/input.go
<<< @../../example/context/date-format/generated/generated.go [generated/generated.go]
:::

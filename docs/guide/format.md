# Input and output formats

Goverter supports three different input output format combinations. This guide
is for you to decide the formats you want to use.

[[toc]]

## interface to struct

This is the default when using
[`goverter:converter`](../reference/converter.md).

**Pros**:

- All goverter features are supported
- The interface is usable before goverter generated the implementation.
  - reduces the occurence of compile errors because of missing or outdated
  generated implementation.
  - allows using generated methods in custom methods

**Cons**:

- You need to initialize the implementation.
- You need to call methods on the struct to execute conversions

::: details Example (click me)
::: code-group
<<< @../../example/format/interfacetostruct/input.go
<<< @../../example/format/common/common.go
<<< @../../example/format/interfacetostruct/generated/generated.go [generated/generated.go]
:::

## variables to assign-variable

This is the default when using
[`goverter:variables`](../reference/converter.md).

**Pros**:

- All goverter features are supported.
- The variables are usable before goverter generated the implementation.
  - reduces the occurence of compile errors because of missing or outdated
    generated implementation.
  - allows using generated functions in custom methods
- You can execute conversions directly without having to initializing a struct

**Cons**:

- Possible runtime overhead when Go cannot optimizen variables the same as
  functions. Benchmark your use-case, if speed is really important to you.

::: details Example (click me)
::: code-group
<<< @../../example/format/assignvariables/input.go
<<< @../../example/format/common/common.go
<<< @../../example/format/assignvariables/input.gen.go
:::

## interface to functions

**Pros**:

- You can execute conversions directly without having to initializing a struct

**Cons**:

- The interface is only used for defining the conversions and is otherwise not
  usable
- The use-case [`map` method with converter](../reference/map.md#method-with-converter) is
  unsupported without replacement.

::: details Example (click me)
::: code-group
<<< @../../example/format/interfacefunction/input.go
<<< @../../example/format/common/common.go
<<< @../../example/format/interfacefunction/generated/generated.go
:::

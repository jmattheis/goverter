# Setting: context

## context ARG

`context ARG` is a can be defined as
[custom function comment](./define-settings.md#custom-function).

`context` defines the `ARG` as context. "Context" are additional arguments that
may be used in other custom functions like [`default`](./default.md),
[`extend`](./extend.md), or
[`map CUSTOM`](./map.md#map-source-path-target-method).

<!-- prettier-ignore -->
::: details Example (click me)
::: code-group
<<< @../../example/context/database/input.go
<<< @../../example/context/database/generated/generated.go [generated/generated.go]
:::

<!-- prettier-ignore -->
::: details Example 2 (click me)
::: code-group
<<< @../../example/context/date-format/input.go
<<< @../../example/context/date-format/generated/generated.go [generated/generated.go]
:::

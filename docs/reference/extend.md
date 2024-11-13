# Setting: extend

`extend [PACKAGE:]FUNC...` can be defined as [CLI argument](./define-settings.md#cli) or
[conversion comment](./define-settings.md#conversion).

If a type cannot be converted automatically, you can manually define an
implementation with `:extend` for the missing mapping. Keep in mind,
Goverter will use the custom implementation every time, when the source and
target type matches.

The `FUNC` may have the signatures described in [Signature:
Optional Source](./signature.md#signature-required-source).

You can optionally define the `PACKAGE` where `FUNC` is located by separating
the `PACKAGE` and `FUNC` with a `:`(colon). If no package is defined, then the
package of the conversion method is used.

`TYPE` can be a regex if you want to include multiple methods from the same
package. E.g. `extend IntTo.*`.

Here are some examples using `extend`.

::: details Simple (click to expand)
::: code-group
<<< @../../example/extend-local/input.go
<<< @../../example/extend-local/generated/generated.go [generated/generated.go]
:::

::: details Complex (click to expand)
::: code-group
<<< @../../example/extend-local-complex/input.go
<<< @../../example/extend-local-complex/generated/generated.go [generated/generated.go]
:::

::: details With `PACKAGE` (click to expand)
::: code-group
<<< @../../example/extend-external/input.go
<<< @../../example/extend-external/generated/generated.go [generated/generated.go]
:::

::: details With Conversion interface (click to expand)
::: code-group
<<< @../../example/extend-local-with-converter/input.go
<<< @../../example/extend-local-with-converter/generated/generated.go [generated/generated.go]
:::

::: details With context (click to expand)
::: code-group
<<< @../../example/context/date-format/input.go
<<< @../../example/context/date-format/generated/generated.go [generated/generated.go]
:::

::: details With error (click to expand)
::: code-group
<<< @../../example/extend-with-error/input.go
<<< @../../example/extend-with-error/generated/generated.go [generated/generated.go]
:::

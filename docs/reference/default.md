# Setting: default

`default [PACKAGE:]FUNC` can be defined as [method
comment](./define-settings.md#method).

By default the target object is initialized with [zero values | Go
docs](https://go.dev/tour/basics/12). With `default` you can instruct Goverter
to use `FUNC` as default target value or constructor for the target value.

The `FUNC` may have the signatures described in [Signature: Optional
Source](./signature.md#signature-optional-source).

You can optionally define the `PACKAGE` where `FUNC` is located by separating
the `PACKAGE` and `FUNC` with a `:`(colon). If no package is defined, then the
package of the conversion method is used.

::: code-group
<<< @../../example/default/input.go
<<< @../../example/default/generated/generated.go [generated/generated.go]
:::

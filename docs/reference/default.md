# default

`default [PACKAGE:]METHOD` can be defined as [method
comment](./define-settings.md#method).

By default the target object is initialized with [zero values | Go
docs](https://go.dev/tour/basics/12). With `default` you can instruct Goverter
to use `METHOD` as default target value or constructor for the target value.

The `METHOD` may be everything supported in [extend
Signatures](./extend.md#signatures) with the addition that the source
parameter is optional.

::: code-group
<<< @../../example/default/input.go
<<< @../../example/default/generated/generated.go [generated/generated.go]
:::

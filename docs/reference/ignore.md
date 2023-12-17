# ignore

`ignore FIELD...` can be defined as [method comment](./define-settings.md#method).

If certain fields shouldn't be converted, are missing on the source struct, or
aren't needed, then you can use `ignore` to ignore these fields.

`ignore` accepts multiple fields separated by spaces. If you want a more global
approach see [ignoreMissing](./ignoreMissing.md) or
[ignoreUnexported](./ignoreUnexported.md)

::: code-group
<<< @../../example/ignore/input.go
<<< @../../example/ignore/generated/generated.go [generated/generated.go]
:::

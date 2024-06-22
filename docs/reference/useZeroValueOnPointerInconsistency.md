# Setting: useZeroValueOnPointerInconsistency

`useZeroValueOnPointerInconsistency [yes,no]` is a
[boolean setting](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

By default, goverter cannot automatically convert `*T` to `T` because it's
unclear how `nil` should be handled. `T` in this example can be any type.

Enable `useZeroValueOnPointerInconsistency` to instruct goverter to use the
zero value of `T` when having the problem above. See [zero values | Go
docs](https://go.dev/tour/basics/12)

::: code-group
<<< @../../example/use-zero-value-on-pointer-inconsistency/input.go
<<< @../../example/use-zero-value-on-pointer-inconsistency/generated/generated.go [generated/generated.go]
:::

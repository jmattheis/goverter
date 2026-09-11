# Setting: useEmptyOnNil

`useEmptyOnNil [yes,no]`, `useEmptyOnNil:slice [yes,no]` and
`useEmptyOnNil:map [yes,no]` are
[boolean settings](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). These settings are
[inheritable](./define-settings.md#inheritance).

By default, goverter converts a `nil` slice/map to `nil`.

* Set `useEmptyOnNil` to enable both of the settings below.
* Set `useEmptyOnNil:map` to convert `nil` maps to empty maps (`map[K]V{}`).
* Set `useEmptyOnNil:slice` to convert `nil` slices to empty slices (`[]T{}`).

::: code-group
<<< @../../example/use-empty-on-nil/input.go
<<< @../../example/use-empty-on-nil/generated/generated.go [generated/generated.go]
:::

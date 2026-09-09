# Setting: useEmptyOnNil

`useEmptyOnNil [yes,no]`, `useEmptyOnNil:slice [yes,no]` and
`useEmptyOnNil:map [yes,no]` are
[boolean settings](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). These settings are
[inheritable](./define-settings.md#inheritance).

By default, goverter passes `nil` slices and maps through: a `nil` source
produces a `nil` target. When encoding to JSON, a `nil` slice is marshalled
as `null` (an empty slice as `[]`), and a `nil` map as `null` (an empty map
as `{}`).

Enable `useEmptyOnNil` to initialize both slices and maps unconditionally
with `make(..., len(source))`, so a `nil` source produces an empty
(non-`nil`) target. `len(nil) == 0`, so this is safe for `nil` input. Use
`useEmptyOnNil:slice` or `useEmptyOnNil:map` to enable it only for slices or
only for maps.

::: code-group
<<< @../../example/use-empty-on-nil/input.go
<<< @../../example/use-empty-on-nil/generated/generated.go [generated/generated.go]
:::

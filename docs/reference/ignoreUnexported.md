# ignoreUnexported

`ignoreUnexported [yes,no]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli), [converter
comment](./define-settings.md#converter) or [method
comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

If a struct has multiple **unexported** fields that should be ignored, then you
can enable `ignoreUnexported`, to ignore these.

::: danger
Using this setting is not recommended, because this can easily lead to
unwanted behavior. When a struct is having unexported fields, you most likely
have to call a custom constructor method to correctly instantiate this type.
:::

::: code-group
<<< @../../example/ignore-unexported/input.go
<<< @../../example/ignore-unexported/generated/generated.go [generated/generated.go]
:::

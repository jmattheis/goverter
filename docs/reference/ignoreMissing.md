# Setting: ignoreMissing

`ignoreMissing [yes,no]` is a [boolean setting](./define-settings.md#boolean) and
can be defined as [CLI argument](./define-settings.md#cli), [conversion
comment](./define-settings.md#conversion) or [method
comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

If the source struct has multiple missing **exported** fields that should be
ignored, then you can enable `ignoreMissing`, to ignore these.

::: danger
Using this setting is not recommended, because this can easily lead to
unwanted behavior when e.g. renaming fields on a struct and forgetting to
change the goverter converter accordingly.
:::

::: code-group
<<< @../../example/ignore-missing/input.go
<<< @../../example/ignore-missing/generated/generated.go [generated/generated.go]
:::

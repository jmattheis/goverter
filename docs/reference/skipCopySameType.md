# Setting: skipCopySameType

`skipCopySameType [yes,no]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli), [conversion
comment](./define-settings.md#conversion) or [method
comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

Goverter deep copies instances when converting the source to the target type.
With `goverter:skipCopySameType` you instruct Goverter to skip copying instances
when the source and target type is the same.

::: code-group
<<< @../../example/skip-copy-same-type/input.go
<<< @../../example/skip-copy-same-type/generated/generated.go [generated/generated.go]
:::

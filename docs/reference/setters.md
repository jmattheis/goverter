# Setting: setters

## setters:enabled yes|no

`setters:enabled [yes|no]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

If the target struct has setter methods, then you can enable setters to use the
setters in the converters.

::: details Example (click me)
::: code-group
<<< @../../example/setters/input.go
<<< @../../example/setters/generated/generated.go [generated/generated.go]
:::

## setters:preferred yes|no

`setters:preferred [yes|no]` is a
[boolean setting](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

If the target struct has a field and setter with the same name, then the
converter prefers to use the setter method instead of the field.

::: details Example (click me)
::: code-group
<<< @../../example/setters-preferred/input.go
<<< @../../example/setters-preferred/generated/generated.go [generated/generated.go]
:::

`setters:regex REGEX` can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance). Default `Set([A-Z].*)`.

`setters:regex` allows you to define a regex to extract a field name from a
setter. The regex must have exactly one capture group. The results are used to
find matching field names and matching getters.

::: details Example (click me)
::: code-group
<<< @../../../example/setters-regex/input.go
<<< @../../../example/setters-regex/generated/generated.go [generated/generated.go]
:::

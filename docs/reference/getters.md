# Setting: getters

## getters:enabled yes|no

getters:enabled [yes|no]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

If the source struct has getter methods, then you can enable getters to implicitly
use getters in the converters.

::: details Example (click me)
::: code-group
<<< @../../example/getters/input.go
<<< @../../example/getters/generated/generated.go [generated/generated.go]
:::

## getters:preferred yes|no

`getters:preferred [yes|no]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

If the source struct has a field and a getter with the same name, then the
converter prefers to use the getter method instead of the field.

::: details Example (click me)
::: code-group
<<< @../../example/getters-preferred/input.go
<<< @../../example/getters-preferred/generated/generated.go [generated/generated.go]
:::

`getters:regex REGEX` can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance). Default `Get{{.}}`.

`getters:regex` allows you to define a template to build a getter name from
a field name. The field name is either the field name directly or extracted
from a setter using [`setters:regex`](./setters.md#gettersregex).

::: details Example (click me)
::: code-group
<<< @../../../example/getters-regex/input.go
<<< @../../../example/getters-regex/generated/generated.go [generated/generated.go]
:::

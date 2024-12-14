# Setting: arg

## arg:context:regex REGEX

`arg:context:regex REGEX` can be defined as [CLI
argument](./define-settings.md#cli), [conversion
comment](./define-settings.md#conversion) or [method
comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance). Default _unset_.

`arg:context:regex` allows you define a regex that automatically defines
arguments as [`context`](./context.md) if the name matches.

::: code-group
<<< @../../example/context/regex/input.go
<<< @../../example/context/regex/generated/generated.go [generated/generated.go]
:::

# Setting: matchIgnoreCase

`matchIgnoreCase [yes,no]` is a
[boolean setting](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

Enable `matchIgnoreCase` to instructs goverter to match fields ignoring
differences in capitalization. If there are multiple matches then goverter will
prefers an exact match (if present) or reports an error. Use
[`map`](./map.md) to fix an ambiquous match error.

::: code-group
<<< @../../example/match-ignore-case/input.go
<<< @../../example/match-ignore-case/generated/generated.go [generated/generated.go]
:::

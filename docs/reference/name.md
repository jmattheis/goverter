# Setting: name

`name NAME` can be defined as [CLI argument](./define-settings.md#cli) or
[conversion comment](./define-settings.md#conversion).

`name` instructs goverter to use the given *name* for the generated struct. By
default goverter will use the interface name and append `Impl` at the end.

::: code-group
<<< @../../example/name-struct/input.go
<<< @../../example/name-struct/generated/generated.go [generated/generated.go]
:::

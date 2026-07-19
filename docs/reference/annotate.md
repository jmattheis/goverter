# Setting: annotate:unmapped

`annotate:unmapped [yes,no]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli), [conversion
comment](./define-settings.md#conversion) or [method
comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

When enabled, goverter adds a comment `// FIELD: SETTING` to the generated
code for each struct field that isn't converted, naming the setting that
caused the field to be unset:

- `// FIELD: ignore` for fields ignored via [`ignore`](./ignore.md)
- `// FIELD: ignoreMissing` for fields without a matching source field via
  [`ignoreMissing`](./ignoreMissing.md)
- `// FIELD: ignoreUnexported` for unexported fields via
  [`ignoreUnexported`](./ignoreUnexported.md)

::: code-group
<<< @../../example/annotate-unmapped/input.go
<<< @../../example/annotate-unmapped/generated/generated.go [generated/generated.go]
:::

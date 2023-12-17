# useUnderlyingTypeMethods

`useUnderlyingTypeMethods [yes,no]` is a
[boolean setting](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[converter comment](./define-settings.md#converter) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

For each type conversion, goverter will check if there is an existing method
that can be used. For named types only the type itself will be checked and not
the underlying type. For this type:

```go
type InputID  int
```

`InputID` would be the _named_ type and `int` the _underlying_ type.

With `useUnderlyingTypeMethods`, goverter will check all named/underlying
combinations.

- named -> underlying
- underlying -> named
- underlying -> underlying

::: code-group
<<< @../../example/use-underlying-type-methods/input.go
<<< @../../example/use-underlying-type-methods/generated/generated.go [generated/generated.go]
:::

Without the setting only the custom method with this signature could be used
for the conversion. `func ConvertUnderlying(s InputID) OutputID`

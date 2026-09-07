# Setting: useEmptyMapOnNil

`useEmptyMapOnNil [yes,no]` is a
[boolean setting](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

By default, goverter passes `nil` maps through: a `nil` source map
produces a `nil` target map. When encoding to JSON, a `nil` map is
marshalled as `null` while an empty map is marshalled as `{}`.

Enable `useEmptyMapOnNil` to instruct goverter to always initialize the
target map with `make(..., len(source))`, so a `nil` source produces an
empty (non-`nil`) map. `len(nil) == 0`, so this is safe for `nil` input.

```go
// goverter:converter
// goverter:useEmptyMapOnNil
type Converter interface {
    Convert(source map[string]Input) map[string]Output
}
```

With the setting enabled, the generated code initializes the target
unconditionally instead of guarding with `if source != nil`:

```go
var outputMap map[string]Output
outputMap = make(map[string]Output, len(source))
for key, value := range source {
    outputMap[key] = Convert(value)
}
return outputMap // nil in -> empty map out -> JSON {}
```

See also [`useEmptySliceOnNil`](./useEmptySliceOnNil.md) for the equivalent
behavior for slices.

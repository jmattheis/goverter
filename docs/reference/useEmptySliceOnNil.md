# Setting: useEmptySliceOnNil

`useEmptySliceOnNil [yes,no]` is a
[boolean setting](./define-settings.md#boolean) and can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

By default, goverter passes `nil` slices through: a `nil` source slice
produces a `nil` target slice. When encoding to JSON, a `nil` slice is
marshalled as `null` while an empty slice is marshalled as `[]`.

Enable `useEmptySliceOnNil` to instruct goverter to always initialize the
target slice with `make(..., len(source))`, so a `nil` source produces an
empty (non-`nil`) slice. `len(nil) == 0`, so this is safe for `nil` input.

```go
// goverter:converter
// goverter:useEmptySliceOnNil
type Converter interface {
    Convert(source []Input) []Output
}
```

With the setting enabled, the generated code initializes the target
unconditionally instead of guarding with `if source != nil`:

```go
var outputList []Output
outputList = make([]Output, len(source))
for i := 0; i < len(source); i++ {
    outputList[i] = Convert(source[i])
}
return outputList // nil in -> [] out -> JSON []
```

See also [`useEmptyMapOnNil`](./useEmptyMapOnNil.md) for the equivalent
behavior for maps.

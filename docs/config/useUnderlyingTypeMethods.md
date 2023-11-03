`useUnderlyingTypeMethods [yes,no]` is a
[boolean setting](config/define.md#boolean) and can be defined as
[CLI argument](config/define.md#cli),
[converter comment](config/define.md#converter) or
[method comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

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

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
// goverter:extend ConvertUnderlying
type Converter interface {
    // goverter:useUnderlyingTypeMethods
    Convert(source Input) Output
}
func ConvertUnderlying(s int) string {
    return ""
}
// these would be used too
// func ConvertUnderlying(s int) OutputID
// func ConvertUnderlying(s InputID) string

type InputID  int
type OutputID string
type Input struct  { ID InputID  }
type Output struct { ID OutputID }
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
    var output example.Output
    output.ID = example.OutputID(example.ConvertUnderlying(int(source.ID)))
    return output
}
```

<!-- tabs:end -->

Without the setting only the custom method with this signature could be used
for the conversion. `func ConvertUnderlying(s InputID) OutputID`

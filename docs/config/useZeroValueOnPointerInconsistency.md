`useZeroValueOnPointerInconsistency [yes,no]` is a
[boolean setting](config/define.md#boolean) and can be defined as
[CLI argument](config/define.md#cli),
[converter comment](config/define.md#converter) or
[method comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

By default, goverter cannot automatically convert `*T` to `T` because it's
unclear how `nil` should be handled. `T` in this example can be any type.

Enable `useZeroValueOnPointerInconsistency` to instruct goverter to use the
zero value of `T` when having the problem above. See [zero values | Go
docs](https://go.dev/tour/basics/12)

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
	// goverter:useZeroValueOnPointerInconsistency
	Convert(source Input) Output
}

type Input struct {
	Name *string
	Age  int
}

type Output struct {
	Name string
	Age  int
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
	var exampleOutput example.Output
	var xstring string
	if source.Name != nil {
		xstring = *source.Name
	}
	exampleOutput.Name = xstring
	exampleOutput.Age = source.Age
	return exampleOutput
}
```

<!-- tabs:end -->

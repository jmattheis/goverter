`ignoreMissing [yes,no]` is a [boolean setting](config/define.md#boolean) and
can be defined as [CLI argument](config/define.md#cli), [converter
comment](config/define.md#converter) or [method
comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

If the source struct has multiple missing **exported** fields that should be
ignored, then you can enable `ignoreMissing`, to ignore these.

!> Using this setting is not recommended, because this can easily lead to
   unwanted behavior when e.g. renaming fields on a struct and forgetting to
   change the goverter converter accordingly.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:ignoreMissing
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    Age int
    Street string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
	var exampleOutput example.Output
	exampleOutput.Name = source.Name
	return exampleOutput
}
```

<!-- tabs:end -->

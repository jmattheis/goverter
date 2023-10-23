`ignoreUnexported [yes,no]` is a [boolean setting](config/define.md#boolean)
and can be defined as [CLI argument](config/define.md#cli), [converter
comment](config/define.md#converter) or [method
comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

If a struct has multiple **unexported** fields that should be ignored, then you
can enable `ignoreUnexported`, to ignore these.

!> Using this setting is not recommended, because this can easily lead to
   unwanted behavior. When a struct is having unexported fields, you most likely
   have to call a custom constructor method to correctly instantiate this type.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:ignoreUnexported
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    // goverter will skip this field
    age int
    street string
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

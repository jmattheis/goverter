`skipCopySameType [yes,no]` is a [boolean setting](config/define.md#boolean)
and can be defined as [CLI argument](config/define.md#cli), [converter
comment](config/define.md#converter) or [method
comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

Goverter deep copies instances when converting the source to the target type.
With `goverter:skipCopySameType` you instruct Goverter to skip copying instances
when the source and target type is the same.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
// goverter:skipCopySameType
type Converter interface {
	Convert(source Input) Output
}

type Input struct {
	Name *string
    ItemCounts map[string]int
}

type Output struct {
	Name *string
    ItemCounts map[string]int
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
	exampleOutput.ItemCounts = source.ItemCounts
	return exampleOutput
}
```

<!-- tabs:end -->

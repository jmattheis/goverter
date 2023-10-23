`matchIgnoreCase [yes,no]` is a
[boolean setting](config/define.md#boolean) and can be defined as
[CLI argument](config/define.md#cli),
[converter comment](config/define.md#converter) or
[method comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

Enable `matchIgnoreCase` to instructs goverter to match fields ignoring
differences in capitalization. If there are multiple matches then goverter will
prefers an exact match (if present) or reports an error. Use
[`map`](config/map.md) to fix an ambiquous match error.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:matchIgnoreCase
    Convert(Input) Output
}

type Input struct {
    Age int
    Fullname string
}
type Output struct {
    Age int
    FULLNAME string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
	var exampleOutput example.Output
	exampleOutput.Age = source.Age
	exampleOutput.FULLNAME = source.Fullname
	return exampleOutput
}
```

<!-- tabs:end -->

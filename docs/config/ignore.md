`ignore FIELD...` can be defined as [method comment](config/define.md#method).

If certain fields shouldn't be converted, are missing on the source struct, or
aren't needed, then you can use `ignore` to ignore these fields.

`ignore` accepts multiple fields separated by spaces. If you want a more global
approach see [ignoreMissing](config/ignoreMissing.md) or
[ignoreUnexported](config/ignoreUnexported.md)


<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:ignore Age
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    Age int
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

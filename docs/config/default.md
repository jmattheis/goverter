`default [PACKAGE:]METHOD` can be defined as [method
comment](config/define.md#method).

By default the target object is initialized with [zero values | Go
docs](https://go.dev/tour/basics/12). With `default` you can instruct Goverter
to use `METHOD` as default target value or constructor for the target value.

The `METHOD` may be everything supported in [extend
Signatures](config/extend.md#signatures) with the addition that the source
parameter is optional.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
	// goverter:default NewOutput
	Convert(*Input) *Output
}

type Input struct {
	Age  int
	Name *string
}
type Output struct {
	Age  int
	Name *string
}

func NewOutput() *Output {
	name := "jmattheis"
	return &Output{
		Age:  42,
		Name: &name,
	}
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source *example.Input) *example.Output {
	pSimpleOutput := example.NewOutput()
	if source != nil {
		var simpleOutput example.Output
		simpleOutput.Age = (*source).Age
		var pString *string
		if (*source).Name != nil {
			xstring := *(*source).Name
			pString = &xstring
		}
		simpleOutput.Name = pString
		pSimpleOutput = &simpleOutput
	}
	return pSimpleOutput
}
```

<!-- tabs:end -->

If you want to add a [Manual Mapping
(`goverter:map`)](/conversion/mapping?id=manual) or [Ignore a field
(`goverter:ignore`)](/conversion/mapping?id=ignore) on a nested type like a map
or slice, then you've to define another converter method.

Example, you've want to `goverter:map` the `NestedInput.LastName` to `NestedOutput.Surname` for this method.

```go
package example

// goverter:converter
type Converter interface {
    Convert([]Input) []Output
}

type Input struct {
    Name string
    Nested NestedInput
}
type NestedInput struct {
    LastName string
}

type Output struct {
    Name string
    Nested NestedOutput
}
type NestedOutput struct {
    Surname string
}
```

You can't do it like this:
```go
// goverter:converter
type Converter interface {
    // goverter:map LastName Surname
    Convert([]Input) []Output
}
```
because Goverter doesn't doesn't apply the mapping to all sub conversions, as
this could cause unexpected behavior. The correct way would be to define
another conversion method like this:

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    Convert([]Input) []Output

    // goverter:map LastName Surname
    ConvertNested(NestedInput) NestedOutput
}

type Input struct {
    Name string
    Nested NestedInput
}
type NestedInput struct {
    LastName string
}

type Output struct {
    Name string
    Nested NestedOutput
}
type NestedOutput struct {
    Surname string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source []example.Input) []example.Output {
	var exampleOutputList []example.Output
	if source != nil {
		exampleOutputList = make([]example.Output, len(source))
		for i := 0; i < len(source); i++ {
			exampleOutputList[i] = c.exampleInputToExampleOutput(source[i])
		}
	}
	return exampleOutputList
}
func (c *ConverterImpl) ConvertNested(source example.NestedInput) example.NestedOutput {
	var exampleNestedOutput example.NestedOutput
	exampleNestedOutput.Surname = source.LastName
	return exampleNestedOutput
}
func (c *ConverterImpl) exampleInputToExampleOutput(source example.Input) example.Output {
	var exampleOutput example.Output
	exampleOutput.Name = source.Name
	exampleOutput.Nested = c.ConvertNested(source.Nested)
	return exampleOutput
}
```

<!-- tabs:end -->

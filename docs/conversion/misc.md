## Rename generated struct

By default, Goverter will name the struct like the interface but with the Impl
suffix, with `goverter:name` you can manually set the name.

<!-- tabs:start -->

#### **input.go**

```go
// goverter:converter
// goverter:name RenamedConverter
type Converter interface {
    Convert(Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
}
```

#### **generated/generated.go**

```go
import example "goverter/example"

type RenamedConverter struct{}

func (c *RenamedConverter) Convert(source example.Input) example.Output {
	var simpleOutput example.Output
	simpleOutput.Name = source.Name
	return simpleOutput
}
```

<!-- tabs:end -->

## Error Wrapping

You can enable error wrapping with the `-wrapErrors` CLI flag or by adding
`goverter:wrapErrors` to the method or converter. Both parameters don't take
any arguments.

Wrapped errors will contain troubleshooting information like the target
struct's field name related to this error.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
// goverter:extend strconv:Atoi
// goverter:wrapErrors
type Converter interface {
    Convert(source Input) (Output, error)
}

type Input struct {
    PostalCode string
}
type Output struct {
    PostalCode int
}
```

#### **generated/generated.go**

```go
package generated

import (
	"fmt"
	example "goverter/example"
	"strconv"
)

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) (example.Output, error) {
	var exampleOutput example.Output
	xint, err := strconv.Atoi(source.PostalCode)
	if err != nil {
		return exampleOutput, fmt.Errorf("error setting field PostalCode: %w", err)
	}
	exampleOutput.PostalCode = xint
	return exampleOutput, nil
}
```

<!-- tabs:end -->

## Use zero value on pointer inconsistency

By default, goverter cannot automatically convert `*T` to `T` because it's
unclear how `nil` should be handled. `T` in this example can be any type.

You can add `goverter:useZeroValueOnPointerInconsistency` to a conversion method
or the converter interface to instruct goverter to use the zero value of `T`
when having the problem above. See [zero values | Go
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

## Skip copy same type

Goverter deep copies instances when converting the source to the target type.
With `goverter:skipCopySameType` you instruct Goverter to skip copying
instances when the source and target type is the same.

The setting can be enabled on both the converter interface and methods.

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

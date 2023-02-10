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
		var errValue example.Output
		return errValue, fmt.Errorf("error setting field PostalCode: %w", err)
	}
	exampleOutput.PostalCode = xint
	return exampleOutput, nil
}
```

<!-- tabs:end -->

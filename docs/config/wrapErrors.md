`wrapErrors [yes,no]` is a
[boolean setting](config/define.md#boolean) and can be defined as
[CLI argument](config/define.md#cli),
[converter comment](config/define.md#converter) or
[method comment](config/define.md#method). This setting is
[inheritable](config/define.md#inheritance).

Enable `wrapErrors` to instruct goverter to wrap errors returned by
[extend](config/extend.md) methods. Wrapped errors will contain troubleshooting
information like the target struct's field name related to this error.

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

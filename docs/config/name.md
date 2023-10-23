`name NAME` can be defined as [CLI argument](config/define.md#cli) or
[converter comment](config/define.md#converter).

`name` instructs goverter to use the given *name* for the generated struct. By
default goverter will use the interface name and append `Impl` at the end.

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

<p align="center">
    <img width="300" src=".github/logo.svg" />
</p>

<h1 align="center">goverter</h1>
<p align="center"><i>golang converter generator</i></p>

## Features

* Automatic conversion of builtin types
  ([`house` example](https://github.com/jmattheis/goverter/blob/main/example/house)), this includes:
    * slices, maps, named types, primitive types, pointers
    * structs with same fields
* Extend parts of the conversion with your own
  func: [`house` example sql.NullString](https://github.com/jmattheis/goverter/blob/main/example/house/input.go#L9)
* Optional return of an error: [`errors` example](https://github.com/jmattheis/goverter/tree/main/example/errors)
* Awesome error
  messages: [mismatch type test](https://github.com/jmattheis/goverter/blob/main/scenario/7_error_nested_mismatch.yml)
* Helper tags like `goverter:map` for converting a struct with same field types but different
  names: [`house` example](https://github.com/jmattheis/goverter/blob/main/example/house/input.go#L13)

## Example

`input.go`
```go
package example

// goverter:converter
type Converter interface {
  Convert(source []Input) []Output
}

type Input struct {
  Name string
  Age int
}
type Output struct {
  Name string
  Age int
}
```
`generated/generted.go`
```go
package generated

import simple "github.com/jmattheis/goverter/example/simple"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source []simple.Input) []simple.Output {
  simpleOutputList := make([]simple.Output, len(source))
  for i := 0; i < len(source); i++ {
    simpleOutputList[i] = c.simpleInputToSimpleOutput(source[i])
  }
  return simpleOutputList
}
func (c *ConverterImpl) simpleInputToSimpleOutput(source simple.Input) simple.Output {
  var simpleOutput simple.Output
  simpleOutput.Name = source.Name
  simpleOutput.Age = source.Age
  return simpleOutput
}
```

*Logo by [MariaLetta](https://github.com/MariaLetta])*

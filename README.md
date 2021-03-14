# go-genconv

## Features

* Automatic conversion of basic types
  ([`house` example](https://github.com/jmattheis/go-genconv/blob/main/example/house)), this includes:
    * slices, maps, named types, primitive types, pointers
    * structs with same fields
* Extend parts of the conversion with your own
  func: [`house` example sql.NullString](https://github.com/jmattheis/go-genconv/blob/main/example/house/input.go#L9)
* Optional return of an error: [`errors` example](https://github.com/jmattheis/go-genconv/tree/main/example/errors)
* Awesome error
  messages: [mismatch type test](https://github.com/jmattheis/go-genconv/blob/main/scenario/7_error_nested_mismatch.yml)
* Helper tags like `genconv:map` for converting a struct with same field types but different
  names: [`house` example](https://github.com/jmattheis/go-genconv/blob/main/example/house/input.go#L13)

## Example

`input.go`
```go
package example

// genconv:converter
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

import simple "github.com/jmattheis/go-genconv/example/simple"

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
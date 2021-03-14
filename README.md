<p align="center">
    <img width="300" src=".github/logo.svg" />
</p>

<h1 align="center">goverter</h1>
<p align="center"><i>a "type-safe Go converters" generator</i></p>
<p align="center">
    <a href="https://github.com/jmattheis/goverter/actions/workflows/build.yml">
        <img alt="Build Status" src="https://github.com/jmattheis/goverter/actions/workflows/build.yml/badge.svg">
    </a>
     <a href="https://codecov.io/gh/jmattheis/goverter">
        <img alt="codecov" src="https://codecov.io/gh/jmattheis/goverter/branch/main/graph/badge.svg">
    </a>
    <a href="https://goreportcard.com/report/github.com/jmattheis/goverter">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/jmattheis/goverter">
    </a>
    <a href="https://github.com/jmattheis/goverter/releases/latest">
        <img alt="latest release" src="https://img.shields.io/github/release/jmattheis/goverter.svg">
    </a>
</p>

goverter is a tool for creating type-safe converters. All you have to
do is create an interface and execute goverter. The project is meant as
alternative to [jinzhu/copier](https://github.com/jinzhu/copier) that doesn't
use reflection.

## Features

- Automatic conversion of builtin types
  ([`house` example](https://github.com/jmattheis/goverter/blob/main/example/house)), this includes:
  - slices, maps, named types, primitive types, pointers
  - structs with same fields
- Extend parts of the conversion with your own
  function: [`house` example sql.NullString](https://github.com/jmattheis/goverter/blob/main/example/house/input.go#L9)
- Optional return of an error: [`errors` example](https://github.com/jmattheis/goverter/tree/main/example/errors)
- Awesome error
  messages: [mismatch type test](https://github.com/jmattheis/goverter/blob/main/scenario/7_error_nested_mismatch.yml)
- No reflection in the generated code

## Usage

1. Create a go modules project if you haven't done so already
   ```bash
   $ go mod init module-name
   ```
1. Add `goverter` as dependency to your project

   ```bash
   $ go get github.com/jmattheis/goverter`
   ```

   Or install the cli globally with:

   ```bash
   $ go install github.com/jmattheis/goverter/cmd/goverter`
   ```

1. Create your converter interface and mark it with a comment containing `goverter:converter`

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

1. Run `goverter`:

   ```
   $ go run github.com/jmattheis/goverter/cmd/goverter module-name-in-full
   # example
   $ go run github.com/jmattheis/goverter/cmd/goverter github.com/jmattheis/goverter/example/simple
   ```

   If you `go install`'ed, then execute it like this:

   ```
   $ goverter module-name-in-full
   ```

1. goverter created a file at `./generated/generated.go`, it may look like this:

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

## Docs

tbd

_Logo by [MariaLetta](https://github.com/MariaLetta)_

`map [SOURCE-PATH] TARGET [| METHOD]` can be defined as [method comment](config/define.md#method).

## map SOURCE-FIELD TARGET 

If the source and target struct have differently named fields, then you can use
`map` to define the mapping.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:map LastName Surname
    Convert(Input) Output
}

type Input struct {
    Age int
    LastName string
}
type Output struct {
    Age int
    Surname string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
	var exampleOutput example.Output
	exampleOutput.Age = source.Age
	exampleOutput.Surname = source.LastName
	return exampleOutput
}
```

<!-- tabs:end -->

### map SOURCE-PATH TARGET

You can access nested properties by separating the field names with `.`:

```go
// goverter:converter
type Converter interface {
    // goverter:map Nested.LastName Surname
    Convert(Input) Output
}
```

<details>
  <summary>Example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:map Nested.LastName Surname
    Convert(Input) Output
}

type Input struct {
    Age int
    Nested NestedInput
}
type NestedInput struct {
    LastName string
}
type Output struct {
    Age int
    Surname string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
	var exampleOutput example.Output
	exampleOutput.Age = source.Age
	exampleOutput.Surname = source.Nested.LastName
	return exampleOutput
}
```

<!-- tabs:end -->

</details>

## map DOT TARGET

`map . TARGET`

When using `.` as source field inside `map`, then you the instruct Goverter to
use the source struct as source for the conversion to the target property. This
is useful, when you've a struct that is the flattened version of another
struct. See [`autoMap`](config/autoMap.md) for the invert operation of this.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
	// goverter:map . Address
	Convert(FlatPerson) Person
}

type FlatPerson struct {
	Name    string
	Age     int
	Street  string
	ZipCode string
}
type Person struct {
	Name    string
	Age     int
	Address Address
}
type Address struct {
	Street  string
	ZipCode string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.FlatPerson) example.Person {
	var examplePerson example.Person
	examplePerson.Name = source.Name
	examplePerson.Age = source.Age
	examplePerson.Address = c.exampleFlatPersonToExampleAddress(source)
	return examplePerson
}
func (c *ConverterImpl) exampleFlatPersonToExampleAddress(source example.FlatPerson) example.Address {
	var exampleAddress example.Address
	exampleAddress.Street = source.Street
	exampleAddress.ZipCode = source.ZipCode
	return exampleAddress
}
```

<!-- tabs:end -->

## map [SOURCE-PATH] TARGET | METHOD

For `[SOURCE-PATH] TARGET` you can use everything that's described above in this document. 

The `METHOD` may be everything supported in [extend Signatures](config/extend.md#signatures)
with the addition that the source parameter is optional.

Similarely to the extend setting you can reference external packages in
`METHOD` by separating the package path and method name by `:`.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:map URL | PrependHTTPS
    // goverter:map . FullName | GetFullName
    // goverter:map Age | DefaultAge
    // goverter:map Value | strconv:Itoa
    Convert(Input) (Output, error)
}

type Input struct{
    URL string

    FirstName string
    LastName  string

    Value int
}
type Output struct{
    URL      string
    FullName string
    Age      int

    Value string
}

func GetFullName(input Input) string {
    return input.FirstName + " " + input.LastName
}

func PrependHTTPS(url string) string {
    return "https://" + url
}

func DefaultAge() int {
    return 42
}
```

#### **generated/generated.go**

```go
package generated

import (
	example "goverter/example"
	"strconv"
)

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) (example.Output, error) {
	var exampleOutput example.Output
	exampleOutput.URL = example.PrependHTTPS(source.URL)
	exampleOutput.FullName = example.GetFullName(source)
	exampleOutput.Age = example.DefaultAge()
	exampleOutput.Value = strconv.Itoa(source.Value)
	return exampleOutput, nil
}
```

<!-- tabs:end -->

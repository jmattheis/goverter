## Method

If a type cannot be converted automatically, you can manually define an
implementation with `goverter:extend` for the missing mapping. Keep in mind,
Goverter will use the custom implementation every time, when the source and
target type matches.

You can pass multiple extend methods to `goverter:extend`, or define the
comment multiple times.

Note: The function name can be a regular expression:

```go
// search for conversion methods that start with SQLStringTo in converter's package
// goverter:extend SQLStringTo.*
// the example below enables strconv.ParseBool method
// goverter:extend strconv:Parse.*
```

<!-- tabs:start -->

#### **input.go**

```go
package example

import "fmt"

// goverter:converter
// goverter:extend IntToString
type Converter interface {
    Convert(Input) Output
}
type Input struct {
    Name string
    Age int
}
type Output struct {
    Name string
    Age string
}

func IntToString(i int) string {
    return fmt.Sprint(i)
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
	exampleOutput.Age = example.IntToString(source.Age)
	return exampleOutput
}
```

<!-- tabs:end -->

<details>
  <summary>Complex example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
// goverter:extend ExtractFriendNames
type Converter interface {
    Convert(source []InputPerson) []OutputPerson
}

type InputPerson struct {
    Name string
    Friends []InputPerson
}
type OutputPerson struct {
    Name string
    Friends []string
}

func ExtractFriendNames(persons []InputPerson) []string {
    var names []string
    for _, person := range persons {
        names = append(names, person.Name)
    }
    return names
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source []example.InputPerson) []example.OutputPerson {
	var exampleOutputPersonList []example.OutputPerson
	if source != nil {
		exampleOutputPersonList = make([]example.OutputPerson, len(source))
		for i := 0; i < len(source); i++ {
			exampleOutputPersonList[i] = c.exampleInputPersonToExampleOutputPerson(source[i])
		}
	}
	return exampleOutputPersonList
}
func (c *ConverterImpl) exampleInputPersonToExampleOutputPerson(source example.InputPerson) example.OutputPerson {
	var exampleOutputPerson example.OutputPerson
	exampleOutputPerson.Name = source.Name
	exampleOutputPerson.Friends = example.ExtractFriendNames(source.Friends)
	return exampleOutputPerson
}
```

<!-- tabs:end -->

</details>

### Method with converter

You can access the generated converter by defining the converter interface as
first parameter.

```go
func IntToString(c Converter, i int) string {
    // Use c.ConvertSomething()
    return fmt.Sprint(i)
}
```

<details>
  <summary>Example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
// goverter:extend ConvertAnimals
type Converter interface {
    Convert(source Input) Output

    // used only in extend method
    ConvertDogs([]Dog) []Animal
    ConvertCats([]Cat) []Animal
}

type Input struct {
    Animals InputAnimals
}
type InputAnimals struct {
    Cats []Cat
    Dogs []Dog
}
type Output struct {
    Animals []Animal
}

type Cat struct { Name string }
type Dog struct { Name string }

type Animal struct { Name string }

func ConvertAnimals(c Converter, input InputAnimals) []Animal {
    dogs := c.ConvertDogs(input.Dogs)
    cats := c.ConvertCats(input.Cats)
    return append(dogs, cats...)
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) example.Output {
	var exampleOutput example.Output
	exampleOutput.Animals = example.ConvertAnimals(c, source.Animals)
	return exampleOutput
}
func (c *ConverterImpl) ConvertCats(source []example.Cat) []example.Animal {
	var exampleAnimalList []example.Animal
	if source != nil {
		exampleAnimalList = make([]example.Animal, len(source))
		for i := 0; i < len(source); i++ {
			exampleAnimalList[i] = c.exampleCatToExampleAnimal(source[i])
		}
	}
	return exampleAnimalList
}
func (c *ConverterImpl) ConvertDogs(source []example.Dog) []example.Animal {
	var exampleAnimalList []example.Animal
	if source != nil {
		exampleAnimalList = make([]example.Animal, len(source))
		for i := 0; i < len(source); i++ {
			exampleAnimalList[i] = c.exampleDogToExampleAnimal(source[i])
		}
	}
	return exampleAnimalList
}
func (c *ConverterImpl) exampleCatToExampleAnimal(source example.Cat) example.Animal {
	var exampleAnimal example.Animal
	exampleAnimal.Name = source.Name
	return exampleAnimal
}
func (c *ConverterImpl) exampleDogToExampleAnimal(source example.Dog) example.Animal {
	var exampleAnimal example.Animal
	exampleAnimal.Name = source.Name
	return exampleAnimal
}
```

<!-- tabs:end -->

</details>

### Method with error

Sometimes, custom conversion may fail, in this case Goverter allows you to
define a second return parameter which must be type

```go
func IntToString(i int) (string, error) {
    return "42", errors.new("oops something went wrong")
}
```

<details>
  <summary>Example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

import "strconv"

// goverter:converter
// goverter:extend StringToInt
type Converter interface {
	Convert(Input) (Output, error)
}

type Input struct{ Value string }
type Output struct{ Value int }

func StringToInt(value string) (int, error) {
	i, err := strconv.Atoi(value)
	return i, err
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) (example.Output, error) {
	var exampleOutput example.Output
	xint, err := example.StringToInt(source.Value)
	if err != nil {
		return exampleOutput, err
	}
	exampleOutput.Value = xint
	return exampleOutput, nil
}
```

<!-- tabs:end -->

</details>

### External Method

You can extend methods from external packages by separating the package path
with `:` from the method.

```go
// goverter:converter
// goverter:extend strconv:Atoi
// goverter:extend github.com/google/uuid:Atoi
type Converter interface {
	Convert(Input) (Output, error)
}
```

<details>
  <summary>Example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
// goverter:extend strconv:Atoi
type Converter interface {
	Convert(Input) (Output, error)
}

type Input struct{ Value string }
type Output struct{ Value int }
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
	xint, err := strconv.Atoi(source.Value)
	if err != nil {
		return exampleOutput, err
	}
	exampleOutput.Value = xint
	return exampleOutput, nil
}
```

<!-- tabs:end -->

</details>

## Mapping Method

Supported in Goverter v0.12.0.

To define a custom conversion method for one specific field, you can use:

```
goverter:map [Mapping] | [Mapping Method]
```

As `[Mapping]` you can use everything that's described in
[Mapping](/conversion/mapping.md).

The `[Mapping Method]` may have:

-   no parameters
-   one parameter passing the source of the conversion
-   two parameters passing the converter interface and the source of the conversion

You can extend methods from external packages by separating the package path
with `:` from the method.

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

### Mapping Method with error

The `[Mapping Method]` can optionally return an error as second return
parameter.

<!-- tabs:start -->

#### **input.go**

```go
package example

import "strconv"

// goverter:converter
type Converter interface {
	// goverter:map NumberString Number | ParseInt
	Convert(Input) (Output, error)
}

type Input struct {
	Name  string
	NumberString string
}
type Output struct {
	Name  string
	Number int
}

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Input) (example.Output, error) {
	var exampleOutput example.Output
	exampleOutput.Name = source.Name
	xint, err := example.ParseInt(source.NumberString)
	if err != nil {
		return exampleOutput, err
	}
	exampleOutput.Number = xint
	return exampleOutput, nil
}
```

<!-- tabs:end -->

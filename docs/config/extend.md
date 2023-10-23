`extend [PACKAGE:]TYPE...` can be defined as [CLI argument](config/define.md#cli) or
[converter comment](config/define.md#converter).

If a type cannot be converted automatically, you can manually define an
implementation with `:extend` for the missing mapping. Keep in mind,
Goverter will use the custom implementation every time, when the source and
target type matches.

See [signatures](#signature) for possible method signatures.


## extend TYPE

`extend TYPE` allows you to reference a method in the local package. E.g:


`TYPE` can be a regex if you want to include multiple methods from the same
package. E.g. `extend IntTo.*`

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

## extend PACKAGE:TYPE

You can extend methods from external packages by separating the package path
with `:` from the method.

`TYPE` can be a regex if you want to include multiple methods from the same
package. E.g. `extend strconv:A.*`


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


## Signatures

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
define a second return parameter which must of type `error`

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

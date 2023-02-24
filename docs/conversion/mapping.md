## Manual

If the source and target struct have differently named fields, then you can use
`goverter:map` to define the mapping manually.

```
goverter:map [SourceField] [TargetField]
```

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

### Manual Nested

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

## Source Object

When using `.` as source field inside `goverter:map`, then you the instruct
Goverter to use the source struct as source for the conversion to the target
property. This is useful, when you've a struct that is the flattened version of
another struct. See also
[`goverter:autoMap`](/conversion/mapping.md?id=auto-map)

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

## Auto Map

You can use `goverter:autoMap` to automatically match fields from a sub struct
to the target struct. This is useful, when you've a struct that is the
flattened version of another struct.

`goverter:autoMap [Path]` accepts one parameter which is a path to a substruct
on the source struct. You can specify nested substructs by separating the
fields with `.`. Example: `Nested.SubStruct`

If there are ambiguities, then goverter will fail with an error.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
	// goverter:autoMap Address
	Convert(Person) FlatPerson
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
type FlatPerson struct {
	Name    string
	Age     int
	Street  string
	ZipCode string
}
```

#### **generated/generated.go**

```go
package generated

import example "goverter/example"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source example.Person) example.FlatPerson {
	var exampleFlatPerson example.FlatPerson
	exampleFlatPerson.Name = source.Name
	exampleFlatPerson.Age = source.Age
	exampleFlatPerson.Street = source.Address.Street
	exampleFlatPerson.ZipCode = source.Address.ZipCode
	return exampleFlatPerson
}
```

<!-- tabs:end -->

## Ignore

If certain fields shouldn't be converted, are missing on the source struct, or
aren't needed, then you can use `goverter:ignore` to ignore these fields.

`goverter:ignore` accepts multiple fields separated by spaces.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:ignore Age
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    Age int
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
	return exampleOutput
}
```

<!-- tabs:end -->

### Ignore All Missing

If a struct has multiple **exported** fields that should be ignored, then you
can use `goverter:ignoreMissing`, to ignore these. The comment can be added to
both the converter interface or to a specific conversion method.

!> Using this setting is not recommended, because this can easily lead to
   unwanted behavior when e.g. renaming fields on a struct and forgetting to
   change the goverter converter accordingly.

<details>
  <summary>Example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:ignoreMissing
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    Age int
    Street string
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
	return exampleOutput
}
```

<!-- tabs:end -->

</details>

### Ignore All Unexported

If a struct has multiple **unexported** fields that should be ignored, then you
can use `goverter:ignoreUnexported`, to ignore these. The comment can be added
to both the converter interface or to a specific conversion method.

!> Using this setting is not recommended, because this can easily lead to
   unwanted behavior. When a struct is having unexported fields, you most likely
   have to call a custom constructor method to correctly instantiate this type.

<details>
  <summary>Example (click to expand)</summary>

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:ignoreUnexported
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    age int
    street string
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
	return exampleOutput
}
```

<!-- tabs:end -->

</details>

## Case-insensitive matching

Goverter will automatically fields if they have the exactly same name. You can
enable case-insensitive field matching with `goverter:matchIgnoreCase`.

Use this tag only when it is extremely unlikely for the source or the target to
have two fields that differ only in casing. E.g.: converting go-jet generated
model to protoc generated struct. If `goverter:matchIgnoreCase` is present and
Goverter detects an ambiquous match, it either prefers an exact match (if
present) or reports an error. Use `goverter:map` to fix an ambiquous match
error.

<!-- tabs:start -->

#### **input.go**

```go
package example

// goverter:converter
type Converter interface {
    // goverter:matchIgnoreCase
    Convert(Input) Output
}

type Input struct {
    Age int
    Fullname string
}
type Output struct {
    Age int
    FULLNAME string
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
	exampleOutput.FULLNAME = source.Fullname
	return exampleOutput
}
```

<!-- tabs:end -->

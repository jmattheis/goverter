`autoMap PATH` can be defined as [method comment](config/define.md#method).

You can use `autoMap` to automatically match fields from a sub struct to the
target struct. This is useful, when you've a struct that is the flattened
version of another struct.

`autoMap PATH` accepts one parameter which is a path to a substruct on the
source struct. You can specify nested substructs by separating the fields with
`.`. Example: `Nested.SubStruct`

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

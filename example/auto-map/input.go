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

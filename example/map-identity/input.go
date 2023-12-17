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

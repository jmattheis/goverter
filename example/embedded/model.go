package example

type Person struct {
	Address
	Name string
}
type Address struct {
	Street  string
	ZipCode string
}
type FlatPerson struct {
	Name       string
	StreetName string
	ZipCode    string
}

// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package generated

import mapidentity "github.com/jmattheis/goverter/example/map-identity"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source mapidentity.FlatPerson) mapidentity.Person {
	var examplePerson mapidentity.Person
	examplePerson.Name = source.Name
	examplePerson.Age = source.Age
	examplePerson.Address = c.exampleFlatPersonToExampleAddress(source)
	return examplePerson
}
func (c *ConverterImpl) exampleFlatPersonToExampleAddress(source mapidentity.FlatPerson) mapidentity.Address {
	var exampleAddress mapidentity.Address
	exampleAddress.Street = source.Street
	exampleAddress.ZipCode = source.ZipCode
	return exampleAddress
}
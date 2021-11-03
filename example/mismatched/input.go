//go:generate go run github.com/jmattheis/goverter/cmd/goverter github.com/jmattheis/goverter/example/mismatched
package mismatched

/*
This example demonstrates how to handle structures that have mismatched
pointers, that for historical reasons have to be maintained as-is.

In this example:
- DBCustomer is DBPerson, *DBAddress
- APICustomer is *APIPerson, APIAddress
*/

// goverter:converter
// goverter:extend ToAPICustomer
type Converter interface {
	Convert(customers DBCustomers) APICustomers
	ToAPIPerson(DBPerson) *APIPerson
	ToAPIAddress(DBAddress) APIAddress
}

func ToAPICustomer(c Converter, customer DBCustomer) APICustomer {
	var result APICustomer
	result.APIPerson = c.ToAPIPerson(customer.DBPerson)
	if customer.DBAddress != nil {
		result.APIAddress = c.ToAPIAddress(*customer.DBAddress)
	}
	return result
}

type DBCustomers []DBCustomer

type DBCustomer struct {
	DBPerson
	*DBAddress
}

type DBPerson struct {
	First, Last string
}

type DBAddress struct {
	Suburb, Postcode string
}

type APICustomers []APICustomer

type APICustomer struct {
	*APIPerson
	APIAddress
}

type APIPerson struct {
	First, Last string
}

type APIAddress struct {
	Suburb, Postcode string
}

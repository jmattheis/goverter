//go:generate go run github.com/jmattheis/goverter/cmd/goverter github.com/jmattheis/goverter/example/mismatched
package mismatched

/*
This example demonstrates how to handle structures that have mismatched
pointers (that for historical reasons cannot be changed).

In this example:
- DBCustomer is DBPerson, *DBAddress
- APICustomer is *APIPerson, APIAddress
- Customers are slices of Customer

In this example, the only conversion that goverter can't infer from directives
is *DBAddress -> APIAddress, so a custom converter has to be written.
*/

// goverter:converter
// goverter:extend AddressOrDefault
type Converter interface {
	Convert(customers DBCustomers) APICustomers

	// goverter:map DBPerson APIPerson
	// goverter:map DBAddress APIAddress
	ToApiCustomer(DBCustomer) APICustomer
	ToAPIAddress(DBAddress) APIAddress
}

func AddressOrDefault(c Converter, address *DBAddress) APIAddress {
	if address == nil {
		// use APIAddress with default values if DBAddress is unset
		return APIAddress{}
	}
	return c.ToAPIAddress(*address)
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

package common

import "database/sql"

type DBHouse struct {
	Address    string
	Apartments []DBApartment
}

type DBApartment struct {
	Position   uint
	Owner      DBPerson
	CoResident []DBPerson
}

type DBPerson struct {
	ID         int
	Name       string
	MiddleName sql.NullString
	Friends    []DBPerson
	Info       Info
}

type APIHouse struct {
	Address    string
	Apartments map[uint]APIApartment
}

type APIApartment struct {
	Position   uint
	Owner      APIPerson
	OwnerName  string
	CoResident []APIPerson
}

type APIPerson struct {
	ID         int
	MiddleName *string
	FirstName  *string
	Friends    []APIPerson
	Info       Info
	Age        int
}

type Info struct {
	Birthplace string
}

package house

import (
	"database/sql"
)

// goverter:converter
// goverter:extend SQLStringToPString
type Converter interface {
	ConvertHouse(source DBHouse) APIHouse
	// goverter:map Name FirstName
	// goverter:ignore Age
	ConvertPerson(source DBPerson) APIPerson
	// goverter:map Owner.Name OwnerName
	ConvertApartment(source DBApartment) APIApartment
}

func SQLStringToPString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

type DBHouse struct {
	Address    string
	Apartments map[int]DBApartment
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
}

type APIHouse struct {
	Address    string
	Apartments map[APIRoomNR]APIApartment
}
type APIRoomNR int

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
	Age        int
}

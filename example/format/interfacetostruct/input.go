package house

import (
	"database/sql"

	"github.com/jmattheis/goverter/example/format/common"
)

// goverter:converter
// goverter:output:format struct
// goverter:extend SQLStringToPString
// goverter:extend ConvertToApartmentMap
type Converter interface {
	ConvertHouse(source common.DBHouse) common.APIHouse
	// goverter:map Owner.Name OwnerName
	ConvertApartment(source common.DBApartment) common.APIApartment
	// goverter:map Name FirstName
	// goverter:ignore Age
	ConvertPerson(source common.DBPerson) common.APIPerson
}

func SQLStringToPString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

func ConvertToApartmentMap(c Converter, source []common.DBApartment) map[uint]common.APIApartment {
	m := make(map[uint]common.APIApartment)
	for _, apartment := range source {
		m[apartment.Position] = c.ConvertApartment(apartment) // !! this is not supported in some formats
	}
	return m
}

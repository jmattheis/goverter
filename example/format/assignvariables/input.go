package house

import (
	"database/sql"

	"github.com/jmattheis/goverter/example/format/common"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend SQLStringToPString
// goverter:extend ConvertToApartmentMap
var (
	ConvertHouse func(source common.DBHouse) common.APIHouse
	// goverter:map Name FirstName
	// goverter:ignore Age
	ConvertPerson func(source common.DBPerson) common.APIPerson
	// goverter:map Owner.Name OwnerName
	ConvertApartment func(source common.DBApartment) common.APIApartment
)

func SQLStringToPString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

func ConvertToApartmentMap(source []common.DBApartment) map[uint]common.APIApartment {
	m := make(map[uint]common.APIApartment)
	for _, apartment := range source {
		m[apartment.Position] = ConvertApartment(apartment) // !! this is not supported in some formats
	}
	return m
}

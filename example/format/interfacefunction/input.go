package house

import (
	"database/sql"

	"github.com/jmattheis/goverter/example/format/common"
)

// goverter:converter
// goverter:output:format function
// goverter:extend SQLStringToPString
type Converter interface {
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

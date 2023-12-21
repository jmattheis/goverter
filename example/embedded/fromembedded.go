package example

// goverter:converter
// goverter:output:file generated/fromembedded.go
type FromConverter interface {
	// goverter:map Address.ZipCode ZipCode
	// goverter:map Address.Street StreetName
	FromEmbedded(Person) FlatPerson
}

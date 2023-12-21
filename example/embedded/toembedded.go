package example

// goverter:converter
// goverter:output:file generated/toembedded.go
type ToConverter interface {
	// goverter:map . Address
	ToEmbedded(FlatPerson) Person

	// goverter:map StreetName Street
	ToEmbeddedAddress(FlatPerson) Address
}

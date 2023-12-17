package example

// goverter:converter
type Converter interface {
	// goverter:map URL | PrependHTTPS
	// goverter:map . FullName | GetFullName
	// goverter:map Age | DefaultAge
	// goverter:map Value | strconv:Itoa
	Convert(Input) (Output, error)
}

type Input struct {
	URL string

	FirstName string
	LastName  string

	Value int
}
type Output struct {
	URL      string
	FullName string
	Age      int

	Value string
}

func GetFullName(input Input) string {
	return input.FirstName + " " + input.LastName
}
func PrependHTTPS(url string) string { return "https://" + url }
func DefaultAge() int                { return 42 }

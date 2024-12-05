package update

// goverter:converter
// goverter:update:ignoreZeroValueField
// ^- this setting can be specified on the conversion method
type WithIgnore interface {
	// goverter:update target
	Convert(source Input, target *Output)
}

// goverter:converter
type WithoutIgnore interface {
	// goverter:update target
	Convert(source Input, target *Output)
}

type Input struct {
	Name    string
	Age     int
	Cool    bool
	Address *string
}

type Output struct {
	Name    string
	Age     int
	Cool    bool
	Address *string
}

package example

// goverter:converter
// goverter:extend ConvertUnderlying
type Converter interface {
	// goverter:useUnderlyingTypeMethods
	Convert(source Input) Output
}

func ConvertUnderlying(s int) string {
	return ""
}

// these would be used too
// func ConvertUnderlying(s int) OutputID
// func ConvertUnderlying(s InputID) string

type InputID int
type OutputID string
type Input struct{ ID InputID }
type Output struct{ ID OutputID }

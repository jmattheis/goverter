package example

// goverter:converter
type Converter interface {
	// goverter:map ID Editable | QueryEditable
	// goverter:context db
	Convert(source PostInput, db Database) (PostOutput, error)
}

// goverter:context db
func QueryEditable(id int, db Database) bool {
	return db.AllowedToEdit(id)
}

type Database interface {
	AllowedToEdit(id int) bool
}

type PostInput struct {
	ID   int
	Body string
}
type PostOutput struct {
	ID       int
	Body     string
	Editable bool
}

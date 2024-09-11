package example

// goverter:converter
type Converter interface {
	// goverter:map ID Editable | QueryEditable
	Convert(source PostInput, ctxDatabase Database) (PostOutput, error)
}

func QueryEditable(id int, ctxDatabase Database) bool {
	return ctxDatabase.AllowedToEdit(id)
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

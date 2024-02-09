package patherr

import "fmt"

type Error struct {
	Path  []any
	Inner error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error at path %s: %s", e.Path, e.Inner)
}

func Wrap(err error, path ...any) error {
	if err, ok := err.(*Error); ok {
		err.Path = append(path, err.Path...)
		return err
	}
	return &Error{Inner: err, Path: path}
}

func Key(v any) any      { return v }
func Index(v int) any    { return v }
func Field(v string) any { return v }

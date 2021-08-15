package builder

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

// Path defines the path inside an error message.
type Path struct {
	Prefix     string
	SourceID   string
	TargetID   string
	SourceType string
	TargetType string
}

const pipe = '|'

// Error defines a conversion error.
type Error struct {
	Path  []*Path
	Cause string
}

// NewError creates an error.
func NewError(cause string) *Error {
	return &Error{Cause: cause, Path: []*Path{}}
}

// Lift appends the path to the error.
func (e *Error) Lift(paths ...*Path) *Error {
	e.Path = append(paths, e.Path...)
	return e
}

// ToString converts the error into a string.
func ToString(err *Error) string {
	if len(err.Path) == 0 {
		panic("oops that shouldn't happen")
	}

	end := 2 + (len(err.Path) * 4) - 1
	sourcePath := (len(err.Path) * 2)
	targetPath := sourcePath + 1

	matrix := make([][]rune, end+1)

	for i := 0; i < len(err.Path); i++ {
		path := err.Path[i]
		padding := int(math.Max(float64(len(path.SourceID)), float64(len(path.TargetID))))

		sourceIdx := i * 2
		if path.SourceType != "" {
			matrix[sourceIdx] = append(matrix[sourceIdx], []rune(strings.Repeat(" ", len(path.Prefix)))...)
			matrix[sourceIdx] = append(matrix[sourceIdx], pipe, ' ')
			matrix[sourceIdx] = append(matrix[sourceIdx], []rune(path.SourceType)...)

			for j := sourceIdx + 1; j < sourcePath; j++ {
				matrix[j] = append(matrix[j], []rune(strings.Repeat(" ", len(path.Prefix)))...)
				matrix[j] = append(matrix[j], pipe)
				matrix[j] = append(matrix[j], []rune(strings.Repeat(" ", padding-1))...)
			}
		} else {
			for j := sourceIdx + 1; j < sourcePath; j++ {
				matrix[j] = append(matrix[j], []rune(strings.Repeat(" ", padding+len(path.Prefix)))...)
			}
		}

		matrix[sourcePath] = append(matrix[sourcePath], []rune(path.Prefix)...)
		matrix[sourcePath] = append(matrix[sourcePath], []rune(path.SourceID)...)
		matrix[sourcePath] = append(matrix[sourcePath], []rune(strings.Repeat(" ", padding-len(path.SourceID)))...)

		targetIdx := end - (i * 2)
		if path.TargetType != "" {
			matrix[targetPath] = append(matrix[targetPath], []rune(path.Prefix)...)
			matrix[targetPath] = append(matrix[targetPath], []rune(path.TargetID)...)
			matrix[targetPath] = append(matrix[targetPath], []rune(strings.Repeat(" ", padding-len(path.TargetID)))...)

			for j := targetIdx - 1; j > targetPath; j-- {
				matrix[j] = append(matrix[j], []rune(strings.Repeat(" ", len(path.Prefix)))...)
				matrix[j] = append(matrix[j], pipe)
				matrix[j] = append(matrix[j], []rune(strings.Repeat(" ", padding-1))...)
			}
			matrix[targetIdx] = append(matrix[targetIdx], []rune(strings.Repeat(" ", len(path.Prefix)))...)
			matrix[targetIdx] = append(matrix[targetIdx], pipe, ' ')
			matrix[targetIdx] = append(matrix[targetIdx], []rune(path.TargetType)...)
		} else {
			for j := targetIdx; j >= targetPath; j-- {
				matrix[j] = append(matrix[j], []rune(strings.Repeat(" ", len(path.Prefix)+padding))...)
			}
		}
	}

	buf := bytes.Buffer{}
	for _, line := range matrix {
		_, _ = fmt.Fprintln(&buf, strings.TrimSpace(string(line)))
	}
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, err.Cause)
	return buf.String()
}

package builder

import (
	"github.com/dave/jennifer/jen"
)

// ErrorWrapper generates code that wraps the errors returned from conversion methods, embedding
// extra details to ease on troubleshooting.
type ErrorWrapper func(errStmt *jen.Statement) *jen.Statement

// NoWrap returns error statement as is.
func NoWrap(errStmt *jen.Statement) *jen.Statement {
	return errStmt
}

// Wrap returns generator that wraps the input error statement with fmt.Errorf.
// Input fmt shall not include the wrapping suffix ": %w", it is appended by the method
// and the input err is appended as the last argument.
func Wrap(fmt string, argStmts ...jen.Code) ErrorWrapper {
	return func(errStmt *jen.Statement) *jen.Statement {
		args := []jen.Code{jen.Lit(fmt + ": %w")}
		args = append(args, argStmts...)
		args = append(args, errStmt)
		return jen.Qual("fmt", "Errorf").Call(args...)
	}
}

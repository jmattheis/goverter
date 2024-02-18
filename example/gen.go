//go:generate go run -cover -covermode=atomic github.com/jmattheis/goverter/cmd/goverter gen ./...
//go:generate go run -cover -covermode=atomic github.com/jmattheis/goverter/cmd/goverter gen -cwd ./wrap-errors-using ./
//go:generate go run -cover -covermode=atomic github.com/jmattheis/goverter/cmd/goverter gen -cwd ./protobuf ./
package example

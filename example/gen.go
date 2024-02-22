//go:generate go run -cover -covermode=atomic github.com/jmattheis/goverter/cmd/goverter gen ./...
//go:generate go run -cover -covermode=atomic github.com/jmattheis/goverter/cmd/goverter gen -cwd ./wrap-errors-using ./
//go:generate go run -cover -covermode=atomic github.com/jmattheis/goverter/cmd/goverter gen -cwd ./protobuf ./
//go:generate go run -C ./enum/transform-custom -cover -covermode=atomic -coverpkg "github.com/jmattheis/goverter/...,goverter/example/..." ./goverter gen ./
package example

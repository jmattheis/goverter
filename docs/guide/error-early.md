# Error early

If Goverter is unable to automatically convert the source to the target type,
then Goverter will error. This is to ensure, that nothing unexpected will
happen.

This is only one example, but you'll get the gist of it.

```go
package example

// goverter:converter
type Converter interface {
    Convert([]Input) []Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    Age int
}
```

In the above example, the `Age` field cannot be automatically converted, because
it does not exist on the input type. When running Goverter it'll fail with the
following error message:

```
Error while creating converter method:
    func (goverter/example.Converter).Convert([]goverter/example.Input) []goverter/example.Output

| []goverter/example.Input
|
|     | goverter/example.Input
|     |
|     |
|     |
source[].???
target[].Age
|     |  |
|     |  | int
|     |
|     | goverter/example.Output
|
| []goverter/example.Output

Cannot match the target field with the source entry: "Age" does not exist.
```

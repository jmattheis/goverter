<script setup>
import { data as libVersion } from './version.data.js'
</script>
# Goverter

goverter is a tool for creating type-safe converters. All you have to
do is create an interface and execute goverter. The project is meant as
alternative to [jinzhu/copier](https://github.com/jinzhu/copier) that doesn't
use reflection.

## Getting Started

1. Ensure your `go version` is 1.16 or above

1. Create a go modules project if you haven't done so already

    ```bash
    $ go mod init module-name
    ```

1. Create your converter interface and mark it with a comment containing
   [`goverter:converter`](https://goverter.jmattheis.de/reference/converter)

    `input.go`

    ```go
    package example

    // goverter:converter
    type Converter interface {
      ConvertItems(source []Input) []Output

      // goverter:ignore Irrelevant
      // goverter:map Nested.AgeInYears Age
      Convert(source Input) Output
    }

    type Input struct {
      Name string
      Nested InputNested
    }
    type InputNested struct {
      AgeInYears int
    }
    type Output struct {
      Name string
      Age int
      Irrelevant bool
    }
    ```

    See [Settings](https://goverter.jmattheis.de/reference/settings) for more information.

1. Run `goverter`:

    ```bash-vue
    $ go run github.com/jmattheis/goverter/cmd/goverter@{{ libVersion }} gen ./
    ```

    It's recommended to use an explicit version instead of `latest`. See
    [Installation](https://goverter.jmattheis.de/guide/install) and
    [CLI](https://goverter.jmattheis.de/reference/cli) for more information.

1. goverter created a file at `./generated/generated.go`, it may look like this:

    ```go
    package generated

    import example "goverter/example"

    type ConverterImpl struct{}

    func (c *ConverterImpl) Convert(source example.Input) example.Output {
        var exampleOutput example.Output
        exampleOutput.Name = source.Name
        exampleOutput.Age = source.Nested.AgeInYears
        return exampleOutput
    }
    func (c *ConverterImpl) ConvertItems(source []example.Input) []example.Output {
        var exampleOutputList []example.Output
        if source != nil {
            exampleOutputList = make([]example.Output, len(source))
            for i := 0; i < len(source); i++ {
                exampleOutputList[i] = c.Convert(source[i])
            }
        }
        return exampleOutputList
    }
    ```

    See [Generation](https://goverter.jmattheis.de/explanation/generation) for more information.

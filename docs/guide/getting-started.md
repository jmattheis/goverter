# Getting started

1. Ensure your `go version` is 1.23 or above

1. Create a go modules project if you haven't done so already

    ```bash
    $ go mod init module-name
    ```

1. [Guide: Install Goverter](./install.md)

1. Create your converter interface and mark it with a comment containing
   [`goverter:converter`](../reference/converter.md)

    ::: code-group
    ```go [input.go]
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
    :::

1. Run Goverter. See [Guide: Install Goverter](./install.md)

1. goverter created a file at `./generated/generated.go`, it may look like this:

    ::: code-group
    ```go [generated/generated.go]
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
    :::

## What's next?

You can look through the the [Structs Guide](./struct.md) and the [Settings
Reference](../reference/settings.md) to find out what's possible with goverter.

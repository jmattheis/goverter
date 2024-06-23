# Setting: extend

`extend [PACKAGE:]TYPE...` can be defined as [CLI argument](./define-settings.md#cli) or
[conversion comment](./define-settings.md#conversion).

If a type cannot be converted automatically, you can manually define an
implementation with `:extend` for the missing mapping. Keep in mind,
Goverter will use the custom implementation every time, when the source and
target type matches.

See [signatures](#signatures) for possible method signatures.


## extend TYPE

`extend TYPE` allows you to reference a method in the local package. E.g:


`TYPE` can be a regex if you want to include multiple methods from the same
package. E.g. `extend IntTo.*`


::: code-group
<<< @../../example/extend-local/input.go
<<< @../../example/extend-local/generated/generated.go [generated/generated.go]
:::


<details>
  <summary>Complex example (click to expand)</summary>

::: code-group
<<< @../../example/extend-local-complex/input.go
<<< @../../example/extend-local-complex/generated/generated.go [generated/generated.go]
:::

</details>

## extend PACKAGE:TYPE

You can extend methods from external packages by separating the package path
with `:` from the method.

`TYPE` can be a regex if you want to include multiple methods from the same
package. E.g. `extend strconv:A.*`

::: code-group
<<< @../../example/extend-external/input.go
<<< @../../example/extend-external/generated/generated.go [generated/generated.go]
:::

## Signatures

### Method with converter

You can access the generated converter by defining the converter interface as
first parameter.

```go
func IntToString(c Converter, i int) string {
    // Use c.ConvertSomething()
    return fmt.Sprint(i)
}
```

<details>
  <summary>Example (click to expand)</summary>

::: code-group
<<< @../../example/extend-local-with-converter/input.go
<<< @../../example/extend-local-with-converter/generated/generated.go [generated/generated.go]
:::

</details>

### Method with error

Sometimes, custom conversion may fail, in this case Goverter allows you to
define a second return parameter which must of type `error`

```go
func IntToString(i int) (string, error) {
    return "42", errors.new("oops something went wrong")
}
```

<details>
  <summary>Example (click to expand)</summary>

::: code-group
<<< @../../example/extend-with-error/input.go
<<< @../../example/extend-with-error/generated/generated.go [generated/generated.go]
:::

</details>

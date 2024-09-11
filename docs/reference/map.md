# Setting: map

`map [SOURCE-PATH] TARGET [| FUNC]` can be defined as [method comment](./define-settings.md#method).

## map SOURCE-FIELD TARGET 

If the source and target struct have differently named fields, then you can use
`map` to define the mapping.

::: code-group
<<< @../../example/map-field/input.go
<<< @../../example/map-field/generated/generated.go [generated/generated.go]
:::

### map SOURCE-PATH TARGET

You can access nested properties by separating the field names with `.`:

```go
// goverter:converter
type Converter interface {
    // goverter:map Nested.LastName Surname
    Convert(Input) Output
}
```

<details>
  <summary>Example (click to expand)</summary>

::: code-group
<<< @../../example/map-path/input.go
<<< @../../example/map-path/generated/generated.go [generated/generated.go]
:::

</details>

## map DOT TARGET

`map . TARGET`

When using `.` as source field inside `map`, then you the instruct Goverter to
use the source struct as source for the conversion to the target property. This
is useful, when you've a struct that is the flattened version of another
struct. See [`autoMap`](./autoMap.md) for the invert operation of this.

::: code-group
<<< @../../example/map-identity/input.go
<<< @../../example/map-identity/generated/generated.go [generated/generated.go]
:::

## map [SOURCE-PATH] TARGET | [PACKAGE:]FUNC

For `[SOURCE-PATH] TARGET` you can use everything that's described above in
this document. The `FUNC` may have the signatures described in [Signature:
Optional Source](./signature.md#signature-optional-source).

You can optionally define the `PACKAGE` where `FUNC` is located by separating
the `PACKAGE` and `FUNC` with a `:`(colon). If no package is defined, then the
package of the conversion method is used.

::: code-group
<<< @../../example/map-custom/input.go
<<< @../../example/map-custom/generated/generated.go [generated/generated.go]
:::

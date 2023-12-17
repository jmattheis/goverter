# map

`map [SOURCE-PATH] TARGET [| METHOD]` can be defined as [method comment](./define-settings.md#method).

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

## map [SOURCE-PATH] TARGET | METHOD

For `[SOURCE-PATH] TARGET` you can use everything that's described above in this document. 

The `METHOD` may be everything supported in [extend Signatures](./extend.md#signatures)
with the addition that the source parameter is optional.

Similarely to the extend setting you can reference external packages in
`METHOD` by separating the package path and method name by `:`.

::: code-group
<<< @../../example/map-custom/input.go
<<< @../../example/map-custom/generated/generated.go [generated/generated.go]
:::

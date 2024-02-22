# Setting: wrapErrorsUsing

`wrapErrorsUsing [PACKAGE]` is a [boolean setting](./define-settings.md#boolean)
and can be defined as [CLI argument](./define-settings.md#cli),
[converter comment](./define-settings.md#converter) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance).

Enable `wrapErrorsUsing` to instruct goverter to wrap errors returned by
[`extend`](./extend.md) and
[`map [SOURCE-PATH] TARGET | METHOD`](./map.md#map-source-path-target-method)
using the implementation specified in `PACKAGE`.

The configured `PACKAGE` is required to have four callable identifiers:

- `func Wrap(error, ELEMENT...) error`: used to wrap conversion errors
- `func Key(any) ELEMENT`: used as `ELEMENT` when the errors occurs for a value of a map
- `func Index(int) ELEMENT`: used as `ELEMENT` when the errors occurs for an item of an array / slice
- `func Field(string) ELEMENT`: used as `ELEMENT` when the errors occurs for field inside a struct

The type of `ELEMENT` can be anything that is applicable to the `Wrap`
signature, it must be consistent between all four methonds. When goverter
returns errors in a conversion function it will `Wrap` the original error and
also provide the target path for the failed conversion.

Goverter creates submethods for code reuse, this means that `Wrap` may be
called multiple times. For your implementation it shouldn't matter if the path
is given in one call or multiple. Both these examples should have the same end result:

```go
err := originalErr
err = Wrap(err, Key("jmattheis"), Index(5), Key("abc"), Field("Name"))
```

```go
err := originalErr
err = Wrap(err, Key("abc"), Field("Name"))
err = Wrap(err, Index(5))
err = Wrap(err, Key("jmattheis"))
```

Here is an example using
[github.com/goverter/patherr](https://github.com/goverter/patherr). This is an
implementation for satisfying the requirements from above. You can use it
directly or as a template for creating your own implementation.

::: code-group 
<<< @../../example/wrap-errors-using/input.go 
<<< @../../example/wrap-errors-using/generated/generated.go [generated/generated.go]
<<< @../../example/wrap-errors-using/go.mod
:::

### Minimal impl

Here is an example minimal implementation of the requirements above:

::: code-group 
<<< @../../example/wrap-errors-using/patherr/patherr.go [patherr/patherr.go]
<<< @../../example/wrap-errors-using/minimal_input.go [input.go]
<<< @../../example/wrap-errors-using/generated/minimal.go [generated/minimal.go]
:::

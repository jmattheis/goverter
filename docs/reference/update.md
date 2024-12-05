# Setting: update

## update ARG

`update ARG` is a can be defined as
[method comment](./define-settings.md#method).

`update` instructs goverter to _update_ an existing instance of a struct passed
via an argument named `ARG`.

Constraints:

- The type of `ARG` must be a pointer to a struct.
- The source type must be a struct or a pointer to a struct.

<!-- prettier-ignore -->
::: details Example (click me)
::: code-group
<<< @../../example/update/input.go
<<< @../../example/update/generated/generated.go [generated/generated.go]
:::

## update:ignoreZeroValueField [yes|no]

`update:ignoreZeroValueField [yes|no]` can be defined as
[CLI argument](./define-settings.md#cli),
[conversion comment](./define-settings.md#conversion) or
[method comment](./define-settings.md#method). This setting is
[inheritable](./define-settings.md#inheritance). Default `no`.

`update:ignoreZeroValueField` instructs goverter to skip mapping fields when the
source field is the [zero value](https://go.dev/tour/basics/12) of the type.

There are three subcategories, `ignoreZeroValueField` is the same as enabling /
disabling all these categories together:

- `update:ignoreZeroValueField:basic [yes|no]`. Ignore zero values for
  [basic types](https://go.dev/tour/basics/11)
- `update:ignoreZeroValueField:struct [yes|no]`. Ignore zero values for
  [structs](https://go.dev/tour/basics/11)
- `update:ignoreZeroValueField:nillable [yes|no]`. Ignore zero values for
  nillable types: channels, maps, functions, interfaces, slices

<!-- prettier-ignore -->
::: details Example (click me)
::: code-group
<<< @../../example/update-ignore-zero/input.go
<<< @../../example/update-ignore-zero/generated/generated.go [generated/generated.go]
:::

`update:ignoreZeroValueField:nillable` does include pointers, but due to the
internal design of goverter (checking for `nil` before dereferencing), it's
ignored by default and cannot be changed currently. This setting still does
affect pointers when [`skipCopySameType`](./skipCopySameType.md) is enabled.

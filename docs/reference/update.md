# Setting: update

`update ARG` is a can bedefined as [method comment](./define-settings.md#method).

`update` instructs goverter to _update_ an existing instance of a struct passed
via an argument named `ARG`.

Constraints:

- The type of `ARG` must be a pointer to a struct.
- The source type must be a struct or a pointer to a struct.

::: code-group
<<< @../../example/update/input.go
<<< @../../example/update/generated/generated.go [generated/generated.go]
:::

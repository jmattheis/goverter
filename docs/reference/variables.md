# Setting: variables

`variables` accepts no arguments and can be defined as [conversion
comment](./define-settings.md#conversion).

`variables` instructs goverter to generate an implementation for the given
variables. You can have multiple variables blocks in one package.

See [output](./output.md) to control the output location/package of the
generated converter.

::: code-group
<<< @../../example/format/assignvariables/input.go
<<< @../../example/format/common/common.go
<<< @../../example/format/assignvariables/input.gen.go
:::

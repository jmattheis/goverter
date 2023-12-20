# Setting: converter

`converter` accepts no arguments and can be defined as [converter
comment](./define-settings.md#converter).

`converter` instructs goverter to generate an implementation for the given
interface. You can have multiple converters in one package.

See [output](./output.md) to control the output location/package of the
generated converter.

::: code-group
<<< @../../example/simple/input.go
<<< @../../example/simple/generated/generated.go [generated/generated.go]
:::

# autoMap

`autoMap PATH` can be defined as [method comment](./define-settings.md#method).

You can use `autoMap` to automatically match fields from a sub struct to the
target struct. This is useful, when you've a struct that is the flattened
version of another struct.

`autoMap PATH` accepts one parameter which is a path to a substruct on the
source struct. You can specify nested substructs by separating the fields with
`.`. Example: `Nested.SubStruct`

If there are ambiguities, then goverter will fail with an error.

::: code-group
<<< @../../example/auto-map/input.go
<<< @../../example/auto-map/generated/generated.go [generated/generated.go]
:::

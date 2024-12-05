# Update an existing instance

When you already have an existing instance to a struct and want to _update_ a
couple of fields from a different type, then you can use
[`update ARG`](../reference/update.md) to instruct goverter to create a
conversion that updates an instance passed as argument.

<!-- prettier-ignore -->
::: code-group
<<< @../../example/update/input.go
<<< @../../example/update/generated/generated.go [generated/generated.go]
:::

## Don't override zero values

If you have fields that aren't pointers, then you may want to enable
[`update:ignoreZeroValueField`](../reference/update.md) to skip overriding
those with the unset values.

<!-- prettier-ignore -->
::: code-group
<<< @../../example/update-ignore-zero/input.go
<<< @../../example/update-ignore-zero/generated/generated.go [generated/generated.go]
:::

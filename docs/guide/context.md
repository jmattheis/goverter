# How to pass context to custom functions

You can pass additional parameters to custom functions by defining them as
[context](../reference/signature.md#categories). This can be done by prefixing
the parameter names with `context` or `ctx` then this guide can give you a
basic understanding on how this works with goverter.

Imaging we want to format a `time.Time` to a `string` but have requirements so
that the date format must be changeable at runtime. You can define the time
format as context, and then use it in a custom
[`extend`](../reference/extend.md) function.


::: code-group
<<< @../../example/context/date-format/input.go
<<< @../../example/context/date-format/generated/generated.go [generated/generated.go]
:::

Similarly, you could supply a database handle and query additional values that
are needed in the `target` but missing in the `source`. E.g. in this example if
a specific post is editable:

::: code-group
<<< @../../example/context/database/input.go
<<< @../../example/context/database/generated/generated.go [generated/generated.go]
:::

# How to convert embedded structs

Goverter sees embedded structs like any other field. This means you can access
them via the type name in [`ignore`](../reference/ignore.md) or
[`map`](../reference/map.md).

Given this model.

::: code-group
<<< @../../example/embedded/model.go
:::

Notice the embedded struct `Address` on `Person`.

When defining the conversion method for `Person` (containing the embedded
struct) to `FlatPerson`, you have to manually define the mapping for `Address`
like this:

::: code-group
<<< @../../example/embedded/fromembedded.go
<<< @../../example/embedded/generated/fromembedded.go [generated/fromembedded.go]
:::

In some cases it's useful to use [`autoMap`](../reference/autoMap.md) if the
embedded struct has equally named fields.

---

For the reverse operation you actually have to define two methods because of
the rules defined in [Guide: Configure Nested](./configure-nested.md). For the
given example it would look like this:

::: code-group
<<< @../../example/embedded/toembedded.go
<<< @../../example/embedded/generated/toembedded.go [generated/toembedded.go]
:::

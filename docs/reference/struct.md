# struct

## struct:comment COMMENT

`struct:comment COMMENT` can be defined as [CLI argument](./define-settings.md#cli)
or [converter comment](./define-settings.md#converter).

`struct:comment` instructs goverter to add a comment line to the generated
struct. It can be configured multiple times to add multiline comments. Prefix
your COMMENT with `//` to force single line comment style.

::: code-group
<<< @../../example/struct-comment/input.go
<<< @../../example/struct-comment/generated/generated.go [generated/generated.go]
:::

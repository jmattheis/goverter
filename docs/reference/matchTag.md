# Setting: matchTag

`matchTag TAGNAME` can be defined as [CLI argument](./define-settings.md#cli) or
[conversion comment](./define-settings.md#conversion).

`matchTag` instructs goverter to use the given TAGNAME to exactly match fields between source and
target fields in structs. This takes precedence over, but will fall back to, standard name matching.
If not specified, no tag matching is performed.

::: code-group
<<< @../../example/matchTag/input.go
<<< @../../example/matchTag/generated/generated.go [generated/generated.go]
:::

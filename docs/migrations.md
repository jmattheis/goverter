## v0.x.x to v1.0.x

### Removed goverter:mapExtend

`goverter:mapExtend FIELD METHOD` can be migrated to `goverter:map . FIELD | METHOD`


### Removed goverter:mapIdentity

`goverter:mapIdentity FIELD` can be migrated to `goverter:map . FIELD`

### CLI

The CLI now requires a command, currently there is only the `gen` generate
command. Furthermore, flags refactored and now have a the same names as the
converter comments.

Old Call

```bash
$ goverter -wrapErrors -ignoreUnexportedFields github.com/jmattheis/goverter/example/simple
```

New Call

```bash
$ goverter gen -g wrapErrors -g ignoreUnexported github.com/jmattheis/goverter/example/simple
```

#### Full flag changes

| Old-Flag                               | New-Flag                                                                         |
| -------------------------------------- | -------------------------------------------------------------------------------- |
| `-wrapErrors`                          | `-g wrapErrors`                                                                  |
| `-matchFieldsIgnoreCase`               | `-g matchIgnoreCase`                                                             |
| `-extend Method1,package/path:Method2` | `-g 'extend Method1 package/path:Method2`'                                       |
| `-extend Method1,package/path:Method2` | `-g 'extend Method1' -g 'extend package/path:Method2'`                           |
| `-ignoreUnexportedFields`              | `-g ignoreUnexported`                                                            |
| `-output PATH`                         | `-g 'output:file FILE'` (NOT RECOMMENDED) See [output](config/output.md)         |
| `-packageName NAME`                    | `-g 'output:package :NAME'` (NOT RECOMMENDED) See [output](config/output.md)     |
| `-packagePath PATH`                    | `-g 'output:package PATH'` (NOT RECOMMENDED) See [output](config/output.md)      |
| `-packagePath PATH -packageName NAME`  | `-g 'output:package PATH:NAME'` (NOT RECOMMENDED) See [output](config/output.md) |

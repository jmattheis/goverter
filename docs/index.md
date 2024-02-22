<script setup>
import { data as libVersion } from './version.data.js'
</script>
# Goverter

goverter is a tool for creating type-safe converters. All you have to
do is create an interface and execute goverter. The project is meant as
alternative to [jinzhu/copier](https://github.com/jinzhu/copier) that doesn't
use reflection.

[Getting Started](./guide/getting-started) ᛫
[Installation](./guide/install) ᛫
[CLI](./reference/cli) ᛫
[Config](./reference/settings)

## Features

- **Fast execution**: No reflection is used at runtime
- Automatically converts builtin types: slices, maps, named types, primitive
  types, pointers, structs with same fields
- [Enum support](./guide/enum)
- [Deep copies](https://en.wikipedia.org/wiki/Object_copying#Deep_copy) per
  default and supports [shallow
  copying](https://en.wikipedia.org/wiki/Object_copying#Shallow_copy)
- **Customizable**: [You can implement custom converter methods](./reference/extend)
- [Clear errors when generating the conversion methods](./guide/error-early) if
  - the target struct has unmapped fields
  - types cannot be converted without losing information

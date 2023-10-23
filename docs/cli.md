The CLI allows you to specify package paths with converter interfaces and
[define global settings](config/define.md#cli) that will be applied to all
converters.

```
$ goverter gen --help
Usage:
  goverter gen [OPTIONS] PACKAGE...

PACKAGE(s):
  Define the import paths goverter will use to search for converter interfaces.
  You can define multiple packages and use the special ... golang pattern to
  select multiple packages. See $ go help packages

OPTIONS:
  -g [value], -global [value]:
          apply settings to all defined converters. For a list of available
          settings see: https://goverter.jmattheis.de/#/config/

  -h, --help:
      display this help page

Examples:
  goverter gen ./example/simple ./example/complex
  goverter gen ./example/...
  goverter gen github.com/jmattheis/goverter/example/simple
  goverter gen -g 'ignoreMissing no' -g 'skipCopySameType' ./simple

Documentation:
  Full documentation is available here: https://goverter.jmattheis.de
```

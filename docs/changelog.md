<script setup>
import GH from './GH.vue';
</script>

# Changelog

## unreleased

- Support for passing context to functions. <GH issue="68" pr="166"/>
  - See [Guide: Pass context to functions](./guide/context.md)
  - See [Reference: Signature](./reference/signature.md)
  - Add [`arg:context:regex`](./reference/arg.md)
- Support updating an existing instance of a struct <GH issue="147" pr="170"/>
  - See [Guide: Update an existing instance](./guide/update-instance.md)
  - Add [`update ARG`](./reference/update.md)
  - Add [`update:ignoreZeroValueField [yes|no]`](./reference/update.md)
- Add [`output:raw CODE`](./reference/output.md#output-raw-code) <GH issue="113" pr="168"/>
- Error on duplicated converter signatures. <GH issue="146" pr="166"/>
- Fix panic when using `chan` in conversion functions. <GH issue="165" pr="167"/>

## v1.5.1

- Prevent duplicated enum cases. <GH issue="150" pr="154"/>
- Enable [`enum`](/reference/enum.md) by default when using
  [`variables`](/reference/variables.md). <GH pr="155"/>

## v1.5.0

- Add two input-output formats. See
  [Guide: Input/Output formats](guide/format.md) for help to choose a format.
  The new formats allow you to call top level functions without the need to
  instantiate a struct to call methods on it. <GH issue="77" pr="149"/>
  - Add [`variables`](reference/variables.md) setting
  - Add [`output:format`](reference/output.md#output-format) setting
- Add line numbers of originating definitions to all error messages.
  <GH pr="149"/>

## v1.4.1

- Error when the settings [`enum`](reference/enum.md) and
  [`useUnderlyingTypeMethods`](reference/useUnderlyingTypeMethods.md) conflict.
  <GH issue="141" pr="142"/>
- Add [CLI](./reference/cli.md) commands `help` and `version`.
  <GH issue="144" pr="145"/>
- Improve generation by assigning variables directly if possible.
  <GH issue="97" pr="148"/> 
  :::details Examples
  ```diff
   func (c *ConverterImpl) Convert(source execution.Input) execution.Output {
       var structsOutput execution.Output
  -    var pString *string
       if source.Nested.Name != nil {
           xstring := *source.Nested.Name
  -        pString = &xstring
  +        structsOutput.Name = &xstring
       }
  -    structsOutput.Name = pString
       return structsOutput
   }
  ```
  ***
  ```diff
   func (c *ConverterImpl) ConvertPToP(source []*int) []*int {
       var pIntList []*int
       if source != nil {
           pIntList = make([]*int, len(source))
           for i := 0; i < len(source); i++ {
  -            var pInt *int
               if source[i] != nil {
                   xint := *source[i]
  -                pInt = &xint
  +                pIntList[i] = &xint
               }
  -            pIntList[i] = pInt
           }
       }
       return pIntList
   }
  ```
  :::

## v1.4.0

- Add [Enum Support](guide/enum.md) <GH issue="61" pr="136"/>. Can be disabled
  via [`enum no`](./guide/enum.md#disable-enum-detection-and-conversion).
- Add [`wrapErrorsUsing`](./reference/wrapErrorsUsing.md) <GH pr="134"/>
- Fix panic with go1.22 <GH issue="135" pr="133"/>
- Require go1.18 for building Goverter <GH pr="133"/>
- Add current working directory `-cwd` option to [CLI](./reference/cli.md) <GH pr="134"/>
- Fix error messages when there is an return error mismatch <GH pr="136"/>
- Fix panic when using type params in [`extend`](./reference/extend),
  [`map`](./reference/map) or [`default`](./reference/default). <GH issue="138" pr="139"/>
- Fix `types.Alias`. See [golang#63223](https://github.com/golang/go/issues/63223)

_internals_:

- Require the examples to be up-to-date via CI <GH pr="136"/>
- Fix file permissions in tests <GH pr="136"/>

## v1.3.2

Change generated directory permissions from `777` -> `755` and generated file
permissions from `777` -> `644`. This only affects newly created files and
directories. Existing files and directories will keep their current
permissions. <GH issue="128" pr="129"/>

## v1.3.1

Fix `nil` map conversion. A `nil` map of will be converted to a `nil` map of the
target type. Previously, the target map was instantiated via `make` with a 0
size. <GH issue="126" pr="127"/>

## v1.3.0

- Fix absolute paths in [`output:file`](reference/output) <GH pr="116"/>
- Error on internally overlapping struct methods <GH issue="114" pr="117"/>
- Scan packages with [build constraints](reference/build-constraint).
  <GH issue="111" pr="118"/>
- [`ignore`](reference/ignore) no longer ignores source fields in combination
  with [`autoMap`](reference/autoMap.md) <GH pr="120"/>
- Error on not existing target fields with field settings <GH pr="120"/>
- Error when using a path as target field <GH pr="120"/>

_internals_:

- Migrate documentation to vitepress <GH pr="118"/>

## v1.2.0

- Support using unexported fields, methods and functions when they are
  accessible from the `output:package` <GH issue="104" pr="107"/>
- Fix ignored field settings for conversions from non pointer to pointer types
  <GH issue="104" pr="107"/>
- Improve error messages for `*T` to `T` conversions <GH issue="105" pr="106"/>
- Error on overlapping internal sub methods <GH issue="105" pr="106"/>

_internals_:

- Execute tests in parallel <GH pr="108"/>

## v1.1.1

Fix a panic when using the `error` type inside the conversion.
<GH issue="102" pr="103"/>

## v1.1.0

- Add [`struct:comment`](reference/struct.md) <GH pr="94"/>
- Add [`useUnderlyingTypeMethods`](reference/useUnderlyingTypeMethods)
  <GH issue="78" pr="95"/>
- Add [`default`](reference/default) <GH issue="93" pr="96"/>
- Allow mapping source methods to the target field. [`map`](reference/map)
  supports method calls as source-path. <GH issue="91" pr="99"/>
- Don't panic on `func()` types <GH pr="99"/>
- Improve error messages <GH pr="98"/>
- add outline of conversion generation [generation](explanation/generation.md)
  <GH pr="99"/>

## v1.0.0

- Major rework of documentation. See [settings](reference/settings)
  <GH pr="92"/>
- Rework of the [CLI](reference/cli.md) <GH pr="92"/>
- Improve handling of boolean flags, allow disabling these settings for single
  methods. See [define](reference/define-settings.md) <GH pr="92"/>
- Add `output:file` `output:pattern` settings: [`output`](reference/output.md)
  <GH pr="92"/>
- Remove deprecated `mapExtend` `mapIdentity`: [migrations](guide/migration.md)
  <GH pr="92"/>
- Refactor internals for upcoming features <GH pr="92"/>
- Remove pkg/errors dependency <GH pr="92"/>
- Improve error messages <GH pr="92"/>

See [migrations](guide/migration.md) for instructions to migrate to this
version. If you have problems with this release please create a ticket in this
project.

## v0.18.0

Add [`skipCopySameType`](reference/skipCopySameType.md), this setting instruct
Goverter to skip copying instances when the source and target type is the same.
<GH issue="86" pr="87"/>

## v0.17.5

Prevent unused variables in generated code when empty structs are used
<GH issue="82" pr="83"/>

## v0.17.4

Fix endless loop when converting nested recursive types.
<GH issue="73" pr="74"/>

## v0.17.3

Fix panic when generating a conversion method containing the type
`map[string]interface{}`. <GH issue="71" pr="72"/>

## v0.17.2

Readd go1.16 support for building goverter, it broke with v0.17.1.
<GH issue="69" pr="70"/>

## v0.17.1

Fix generation of types with generic arguments <GH issue="66" pr="67"/>

## v0.17.0

- Add
  [`goverter:useZeroValueOnPointerInconsistency`](reference/useZeroValueOnPointerInconsistency.md)
  <GH issue="64" pr="65"/>
- Allow defining [`matchIgnoreCase`](reference/matchIgnoreCase.md) on the
  converter interface <GH commit="54d5514eddb90c46aacf48f0b3bde609d3d9f1ec"/>

## v0.16.0

- Add [`ignoreMissing`](reference/ignoreMissing.md) <GH pr="60"/>
- Add [`ignoreUnexported`](reference/ignoreUnexported) <GH pr="60"/>

## v0.15.0

- Improve code generation for errors <GH issue="19" pr="57"/>
- Add [`goverter:autoMap`](reference/autoMap.md) <GH issue="55" pr="59"/>

## v0.14.0

- Prevent value copying of source struct pointers if possible. This should fix
  "go vet copylocks" warnings, because some structs should not be copied. See
  <GH issue="39" pr="56"/>
- Due to the change above, the generated code will look different, because
  goverter now splits internal converter methods differently. The overall
  behavior of the implementation shouldn't change. <GH pr="56"/>
- Error on overlapping field mappings. This doesn't change how config is
  evaluated, but it does now error when field mapping config like `goverter:map`
  is at the wrong converter method and would be ignored. <GH pr="56"/>

## v0.13.0

- Fix docs links in error messages
- Allow using external packages in [`map`](reference/map.md)
  <GH issue="47" pr="50"/>
- Allow returning an error in [`map`](reference/map.md) <GH issue="43" pr="50"/>
- Error on misconfiguration <GH issue="8" pr="51"/>

## v0.12.0

- **Deprecation**: `goverter:mapExtend` will be removed soon. See
  [migration](guide/migration.md).
- **Deprecation**: `goverter:mapIdentity` will be removed soon. See
  [migration](guide/migration.md)
- Add support for unnamed structs
  <GH issue="41" commit="6287f10e5e3ae1d8ebcaa50138898242c659aa2d"/>
- Add documentation https://goverter.jmattheis.de/ <GH pr="49"/>
- Allow defining a custom conversion for only a specific field. See
  [`map`](reference/map.md) <GH pr="49"/>

## v0.11.1

Fix pointer slice conversion
<GH issue="40" commit="be38743c1f0e2003bd8fbb54a8df0f35c96c05e4"/>

## v0.11.0

Allow passing self in `mapExtend` method
<GH issue="38" commit="eb133772c13348174561945d39d3491c34403811"/>

## v0.10.1

- Improve error message on compile errors <GH pr="35"/>
- Fix `goverter:mapIdentity` with pointer types
  <GH issue="36" commit="feeafb6120fb6e2873830c3df6df05124557b84b"/>

## v0.10.0

By default, goverter will fail if you don't [`ignore`](reference/ignore.md) all
unexported fields. To automatically ignore all unexported fields, you can enable
[`ignoreUnexported`](./reference/ignoreUnexported.md) <GH issue="31" pr="32"/>

## v0.9.0

Add [`wrapErrors`](reference/wrapErrors.md) <GH pr="29"/>

## v0.8.1

Use extend method when converter method with same method exists
<GH issue="26" commit="12cc6475bd1f27296b06d2b1050a32cec35c81a9"/>

## v0.8.0

Add `mapExtend FIELD METHOD` <GH pr="25"/>

## v0.7.0

Allow setting the packagePath of the generated converter to prevent import
loops. <GH pr="22"/>

## v0.6.3

- Fix generation with generics <GH pr="20"/>
- Exit with non-zero on error <GH pr="21"/>

## v0.6.2

Fix compile errors in code generation with error in return type

## v0.6.1

Error using [`extend`](reference/extend.md) on unexported method <GH pr="17"/>

## v0.6.0

Add [`matchIgnoreCase`](reference/matchIgnoreCase.md) <GH pr="16"/>

## v0.5.0

Add working directory setting. <GH pr="15"/>

## v0.4.0

- Allow external packages in [`extend`](reference/extend.md) <GH pr="14"/>
- Allow regex in function names in [`extend`](reference/extend.md) <GH pr="14"/>

## v0.3.0

Add `mapIdentity FIELD`. <GH issue="12" pr="13"/>

## v0.2.0

- Support nesting in [`map`](reference/map.md). <GH issue="3" pr="5"/>
- Fail on structs with unexported fields. <GH pr="5"/>

## v0.1.2

Fix [`map`](reference/map.md) & [`ignore`](reference/ignore.md) for struct
pointer

## v0.1.1

Add tests

## v0.1.0

Initial Release

# Settings

Before configuring settings here it is useful to understand how [converter
generation works](../explanation/generation.md) and how to [configure nested
settings](../guide/configure-nested.md).

## Converter

These settings can only be defined as [CLI argument](./define-settings.md#cli) or
[converter comment](./define-settings.md#converter).

- [`converter` marker comment for conversion interfaces](./converter.md)
- [`enum [yes|no]` enable / disable enum support](./enum.md#enum-detect)
- [`enum:exclude [PACKAGE:]NAME` exclude wrongly detected enums](./enum.md#enum-exclude)
- [`extend [PACKAGE:]TYPE...` add custom methods for conversions](./extend.md)
- [`name NAME` rename generated struct](./name.md)
- [`output:file FILE` set the output directory for a converter](./output.md#outputfile)
- [`output:package [PACKAGE:]NAME` set the output package for a converter](./output.md#outputpackage)
- [`struct:comment COMMENT` add comments to generated struct](./struct.md#structcomment-comment)

## Method:

These settings can only be defined as [method comment](./define-settings.md#method).

- [`autoMap PATH` automatically match fields from a sub struct to the target struct](./autoMap.md)
- [`default [PACKAGE:]TYPE` define default target value](./default.md)
- [`enum:map SOURCE TARGET` define an enum value mapping](./enum.md#enum-map-source-target)
- [`enum:transform ID CONFIG` use an enum value transformer](./enum.md#enum-transform-id-config)
- [`ignore FIELD...` ignore fields for a struct](./ignore.md)
- [`map [SOURCE-PATH] TARGET [| METHOD]` struct mappings](./map.md)
  - [`map SOURCE-FIELD TARGET` define a field mapping](./map.md#map-source-field-target)
  - [`map SOURCE-PATH TARGET` define a nested field mapping](./map.md#map-source-path-target)
  - [`map . TARGET` map the source type to the target field](./map.md#map-dot-target)
  - [`map [SOURCE-PATH] TARGET| METHOD` map the SOURCE-PATH to the TARGET field by
    using METHOD](./map.md#map-source-path-target-method)


### Method (inheritable)

These settings can be defined as [CLI argument](./define-settings.md#cli),
[converter comment](./define-settings.md#converter) or
[method comment](./define-settings.md#method) and are
[inheritable](./define-settings.md#inheritance).

- [`enum:unknown ACTION|KEY` handle unexpected enum values](./enum.md#enum-default-action)
- [`ignoreMissing [yes,no]` ignore missing struct fields](./ignoreMissing.md) 
- [`ignoreUnexported [yes,no]` ignore unexported struct fields](./ignoreUnexported.md)
- [`matchIgnoreCase [yes,no]` case-insensitive field matching](./matchIgnoreCase.md)
- [`skipCopySameType [yes,no]` skip copying types when the source and target type are the same](./skipCopySameType.md)
- [`useUnderlyingTypeMethods [yes|no]` use underlying types when looking for existing methods](./useUnderlyingTypeMethods.md)
- [`useZeroValueOnPointerInconsistency [yes|no]` Use zero values for `*S` to `T` conversions](./useZeroValueOnPointerInconsistency.md)
- [`wrapErrorsUsing [PACKAGE]` wrap errors using a custom implementation](./wrapErrorsUsing.md)
- [`wrapErrors [yes,no]` wrap errors with extra information](./wrapErrors.md)

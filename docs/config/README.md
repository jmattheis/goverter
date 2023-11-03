Before configuring settings here it is useful to understand how [converter
generation works](generation.md) and how to [configure nested
settings](config/nested.md).

### Converter

These settings can only be defined as [CLI argument](config/define.md#cli) or
[converter comment](config/define.md#converter).

- [`converter` marker comment for conversion interfaces](config/converter.md)
- [`extend [PACKAGE:]TYPE...` add custom methods for conversions](config/extend.md)
- [`name NAME` rename generated struct](config/name.md)
- [`output:file FILE` set the output directory for a converter](config/output.md#file)
- [`output:package [PACKAGE:]NAME` set the output package for a converter](config/output.md#package)
- [`struct:comment COMMENT` add comments to generated struct](config/struct.md#structcomment-comment)

### Method:

These settings can only be defined as [method comment](config/define.md#method).

- [`autoMap PATH` automatically match fields from a sub struct to the target struct](config/autoMap.md)
- [`ignore FIELD...` ignore fields for a struct](config/ignore.md)
- [`map [SOURCE-PATH] TARGET [| METHOD]` struct mappings](config/map.md)
  - [`map SOURCE-FIELD TARGET` define a field mapping](config/map.md#map-source-field-target)
  - [`map SOURCE-PATH TARGET` define a nested field mapping](config/map.md#map-source-path-target)
  - [`map . TARGET` map the source type to the target field](config/map.md#map-dot-target)
  - [`map [SOURCE-PATH] TARGET| METHOD` map the SOURCE-PATH to the TARGET field by
    using METHOD](config/map.md#map-source-path-target-method)


#### Method (inheritable)

These settings can be defined as [CLI argument](config/define.md#cli),
[converter comment](config/define.md#converter) or
[method comment](config/define.md#method) and are
[inheritable](config/define.md#inheritance).

- [`ignoreMissing [yes,no]` ignore missing struct fields](config/ignoreMissing.md) 
- [`ignoreUnexported [yes,no]` ignore unexported struct fields](config/ignoreUnexported.md)
- [`matchIgnoreCase [yes,no]` case-insensitive field matching](config/matchIgnoreCase.md)
- [`skipCopySameType [yes,no]` skip copying types when the source and target type are the same](config/skipCopySameType.md)
- [`useUnderlyingTypeMethods [yes|no]` use underlying types when looking for existing methods](config/useUnderlyingTypeMethods.md)
- [`useZeroValueOnPointerInconsistency [yes|no]` Use zero values for `*S` to `T` conversions](config/useZeroValueOnPointerInconsistency.md)
- [`wrapErrors [yes,no]` wrap errors with extra information](config/wrapErrors.md)

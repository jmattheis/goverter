# Generation

Here is a broad non exhaustive outline what goverter does. Goverter tries to
simplify until there are primitive types to then copy by value.

Definition:

- `CS` current source type
- `CT` current target type
- `NS` next source type
- `NT` next target type

`generate(CS, CT)`:

1. if [`extend`](../reference/extend.md) method exists `func(CS) CT`:
   - use it
   - method generation finished
2. if [`useUnderlyingTypeMethods`](../reference/useUnderlyingTypeMethods.md) is
   enabled:
   - and if `extend` method exists for a underlying type
     - use it
     - method generation finished
3. if [`skipCopySameType`](../reference/skipCopySameType.md) is enabled and
   `CS==CT`:
   - use source value
   - method generation finished
4. `CS is primitive` and `CT is primitive` and `CS=CT`
   - copy source value
   - method generation finished
5. `CS is *NS` and `CT is *NT`:
   - check for `nil` and `generate(NS, NT)`
6. `CS` is not a pointer and `CT is *NT`:
   - get pointer of `generate(CS, NT)`
7. if `CS is *NS` and `CT` is not pointer and
   [`useZeroValueOnPointerInconsistency`](../reference/useZeroValueOnPointerInconsistency.md)
   is enabled:
   - if `CS` is `nil`: use zero value of `CT` otherwise: `generate(NS, CT)`
8. `CS is []NS` and `CT is []NT`
   - iterate over the slice
     - convert slice item: `generate(NS, NT)`
9. `CS is map[Key-NS]Value-NS` and `CT is map[Key-NS]Value-NT`
   - iterate over the map
     - convert the key: `generate(Key-NS, KEY-NT)`
     - convert the value: `generate(Value-NS, Value-NT)`
10. `CS is struct` and `CT is struct`:
    - for each TargetField(TF) in CT:
      - if `TF` is [`ignore`](../reference/ignore.md)d
        - skip
      - if `TF` is unaccessible and
        [`ignoreUnexported`](../reference/ignoreUnexported.md) is not enabled:
        - error: cannot use unexported types
      - get SourceField(SF)
        - from [`map`](../reference/map.md) if defined
        - otherwise:
          - try to find field with same name from `CS`
          - try to find field in paths defined in
            [`autoMap`](../reference/autoMap.md)
          - field matching is case insensitive if
            [`matchIgnoreCase`](../reference/matchIgnoreCase.md) is enabled)
      - if `SF` is a method on `CS`:
        - call `SF`
      - if
        [mapping method MF](../reference/map.md#map-source-path-target-method)
        is defined
        - execute `MF(SF) MappingTarget`
        - ensure `MappingTarget` == `typeof TF`
      - else `generate(SF) TF`
11. error: `CS` cannot be automatically converted to `CT`.

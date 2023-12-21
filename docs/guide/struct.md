# Struct conversion basics

[[toc]]

## How to flatten and unflatten structs

- For unflattening: [`map DOT TARGET`](../reference/map#map-dot-target)
- For flattening: [`autoMap`](../reference/autoMap.md)

## How to map field to a constant

See [`map [SOURCE-PATH] TARGET | METHOD`](../reference/map#map-source-path-target-method) (`DefaultAge`)

## How to ignore errors for missed fields

See [`ignoreMissing`](../reference/ignoreMissing.md)

## How to convert to protobuf generated structs

Here is an example with protobuf generated code:

::: details Example (click me)
::: code-group
<<< @../../example/protobuf/input.go
<<< @../../example/protobuf/event.proto
<<< @../../example/protobuf/pb/event.pb.go [pb/event.pb.go]
<<< @../../example/protobuf/generated/generated.go [generated/generated.go]
:::

## How to map differently named fields

See [`map`](../reference/map.md)

## How to skip conversion of a field

See [`ignore`](../reference/ignore.md)

## How to handle unexported fields

See [Guide: How to handle unexported fields](./unexported-field.md)

## How to convert nested structs

See [Guide: Configure Nested](./configure-nested.md)

## How to convert embedded structs

See [Guide: Structs Embedded](./embedded-structs.md)

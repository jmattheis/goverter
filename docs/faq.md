# Frequently Asked Questions

[[toc]]

## How to convert structs

See [Guide: Struct](./guide/struct.md)

## TypeMismatch: Cannot convert interface{} to interface{}

See below

## TypeMismatch: Cannot convert any to any

Goverter doesn't know how to convert `any` to `any` or `interface{}` to
`interface{}`, you need to define the mapping yourself. If you only want to
pass the value without any conversion you can define it like this:

::: details Example (click me)
::: code-group
<<< @../../example/any-to-any/input.go
<<< @../../example/any-to-any/generated/generated.go [generated/generated.go]
:::

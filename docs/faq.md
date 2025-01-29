# Frequently Asked Questions

[[toc]]

## How to convert structs

See [Guide: Struct](./guide/struct.md)

## Skip conversion (aka shallow copy), e.g for time.Time

Many structs like `time.Time` or `net/url.URL` are best to
just copy instead of convert field-by-field.

To resolve errors like these:

```console
$ goverter gen ./...
Error while creating converter method:
    /home/myuser/code/my-project/converters.go:9
    func (my-project/main.Converter).ConvertFoo(source my-project/main.Source) my-project/main.Target
        [source] my-project/main.Source
        [target] my-project/main.Target

| my-project/main.Source
|
|      | time.Time
|      |
source.CreatedAt.???
target.CreatedAt.wall
|      |         |
|      |         | uint64
|      |
|      | time.Time
|
| my-project/main.Target

Cannot set value for unexported field "wall".

See https://goverter.jmattheis.de/guide/unexported-field
```

You need to add a custom converter, such as:

::: details Example (click me)
::: code-group
<<< @../../example/time/input.go
<<< @../../example/time/generated/generated.go [generated/generated.go]
:::

See [Guide: Structs / Unexported field](./guide/unexported-field.md)

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

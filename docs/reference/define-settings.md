# Define Settings

Settings are can be defined in the following locations:

## CLI

Settings defined globally via the CLI are applied to all Converters seen by
Goverter.

```
goverter -g 'setting yes' -g `setting2 no` ...
```

or

```
goverter -global 'setting yes' -global `setting2 no` ...
```

given this example:

```go
// goverter:converter
type Converter interface {
    Convert(source Input) Output
}
```

the resolved settings would be the same as with the example below for Converter.

## Conversion

When using [`goverter:converter`](converter.md), then you can define settings
on the converter interface by prefixing them with `goverter:`.

```go
// goverter:converter
// goverter:setting yes
// goverter:setting2 no
type Converter interface {
    Convert(source Input) Output
}
```

When using [`goverter:variables`](variables.md), then you can define settings
on the `var`-block by prefixing them with `goverter:`.

```go
// goverter:variables
// goverter:setting yes
// goverter:setting2 no
var (
    Convert func(source Input) Output
)
```

## Method

When using [`goverter:converter`](converter.md), then you can define settings
on the converter methods by prefixing them with `goverter:`.

```go
// goverter:converter
type Converter interface {
    // goverter:setting yes
    // goverter:setting2 no
    Convert(source Input) Output
}
```

When using [`goverter:variables`](variables.md), then you can define settings
on the function variables by prefixing them with `goverter:`.

```go
// goverter:variables
var (
    // goverter:setting yes
    // goverter:setting2 no
    Convert func(source Input) Output
)
```

### Inheritance

Method settings can be inherited for all methods if they are defined on the CLI
or Converter interface. Settings defined on methods take precedence over
inherited settings. So you can enable a feature globally and disable it for one
specific method.

In the example below `feature1` would have the value `yes` for both `Convert1` &
`Convert2` and the value `no` for `Convert3`.

```go
// goverter:converter
// goverter:feature1 yes
type Converter interface {
    Convert1(source Input1) Output1
    Convert2(source Input2) Output2
    // goverter:feature1 no
    Convert3(source Input3) Output3
}
```

## Boolean

Boolean settings have two states: `yes` and `no`

You can enable a setting by defining it without a value:

```
// goverter:setting
```

Or with the `yes` value

```
// goverter:setting yes
```

You can disable a setting with the `no` value:

```
// goverter:setting no
```

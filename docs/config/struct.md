## struct:comment COMMENT

`struct:comment COMMENT` can be defined as [CLI argument](config/define.md#cli)
or [converter comment](config/define.md#converter).

`struct:comment` instructs goverter to add a comment line to the generated
struct. It can be configured multiple times to add multiline comments. Prefix
your COMMENT with `//` to force single line comment style.

<!-- tabs:start -->

#### **input.go**

```go
// goverter:converter
// goverter:struct:comment // MyConverterImpl
// goverter:struct:comment //
// goverter:struct:comment // More detailed
type MultipleSingleLine interface {
    Convert(Input) Output
}

// goverter:converter
// goverter:struct:comment single comment
type SingleComment interface {
    Convert(Input) Output
}

// goverter:converter
// goverter:struct:comment MyConverterImpl
// goverter:struct:comment
// goverter:struct:comment More detailed
type MultiLine interface {
    Convert(Input) Output
}

type Input struct { Name    string }
type Output struct { Name    string }
```

#### **generated/generated.go**

```go
import example "goverter/example"

/*
MyConverterImpl

More detailed
*/
type MultiLineImpl struct{}

// MyConverterImpl
//
// More detailed
type MultipleSingleLineImpl struct{}

// single comment
type SingleCommentImpl struct{}
```

<!-- tabs:end -->

input:
    input.go: |
        package example

        // goverter:converter
        // goverter:extend ConvertUnderlying
        type Converter interface {
            // goverter:useUnderlyingTypeMethods
            Convert(SqlColor) Color
        }

        func ConvertUnderlying(s string) string {
            return ""
        }

        type SqlColor string
        const SqlColorDefault SqlColor = "default"

        type Color string
        const ColorDefault Color = "default"
error: |-
    Error while creating converter method:
        func (github.com/jmattheis/goverter/execution.Converter).Convert(github.com/jmattheis/goverter/execution.SqlColor) github.com/jmattheis/goverter/execution.Color

    | github.com/jmattheis/goverter/execution.SqlColor
    |
    source
    target
    |
    | github.com/jmattheis/goverter/execution.Color

    The conversion between the types
        github.com/jmattheis/goverter/execution.SqlColor
        github.com/jmattheis/goverter/execution.Color

    does qualify for enum conversion but also match an extend method via useUnderlyingTypeMethods.
    You have to disable enum or useUnderlyingTypeMethods to resolve the setting conflict.

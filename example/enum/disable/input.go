package example

import (
	"time"
)

// goverter:converter
// goverter:enum no
type Converter interface {
	Convert(MyDuration) time.Duration
}

type MyDuration int64

const (
	Nanoseconds  MyDuration = 1
	Microseconds MyDuration = 1000 * Nanoseconds
)

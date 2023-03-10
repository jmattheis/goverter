package builder

type ConversionFlag int

const (
	FlagWrapErrors ConversionFlag = iota + 1
	FlagIgnoreMissing
	FlagMatchIgnoreCase
	FlagIgnoreUnexported
	FlagZeroValueOnPtrInconsistency
)

type ConversionFlags map[ConversionFlag]bool

func (c ConversionFlags) Has(flag ConversionFlag) bool {
	return c[flag]
}

func (c ConversionFlags) Set(flag ConversionFlag) {
	c[flag] = true
}

func (c ConversionFlags) Add(flags ConversionFlags) ConversionFlags {
	result := ConversionFlags{}
	for flag, ok := range c {
		if ok {
			result.Set(flag)
		}
	}
	for flag, ok := range flags {
		if ok {
			result.Set(flag)
		}
	}
	return result
}

package flagutil

const (
	boolFlagName  = "bool_flag"
	testFlagName  = "string_flag"
	testFlagValue = "the-value"
)

type flagPackage int

const (
	fpNone flagPackage = iota
	fpBase
	fpSPF13
)

func (p flagPackage) String() string {
	var s string
	switch p {
	case fpNone:
		s = "no-flags"
	case fpBase:
		s = "base"
	case fpSPF13:
		s = "spf13"
	}
	return s
}

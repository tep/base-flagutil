package flagutil

import (
	"flag"
	"github.com/spf13/pflag"
)

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

func resetForTest(name string) {
	merged = make(map[string]bool)
	cmdlineArgs = []string{}

	flag.CommandLine = flag.NewFlagSet("test:"+name, flag.ContinueOnError)
	pflag.CommandLine = pflag.NewFlagSet("test:"+name, pflag.ContinueOnError)
}

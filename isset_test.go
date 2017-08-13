package flagutil

import (
	"flag"
	"testing"

	"github.com/spf13/pflag"
)

type isSetTestcase struct {
	name string
	fpkg flagPackage
	want string
	args []string
}

func newIsSetTestcase(fpkg flagPackage, set bool) *isSetTestcase {
	var name string

	switch {
	case fpkg == fpNone:
		name = "no-flags"

	case fpkg == fpBase && set:
		name = "base-isset"

	case fpkg == fpBase:
		name = "base-noset"

	case fpkg == fpSPF13 && set:
		name = "spf13-isset"

	case fpkg == fpSPF13:
		name = "spf13-noset"
	}

	tc := &isSetTestcase{
		name: name,
		fpkg: fpkg,
	}

	if fpkg == fpNone {
		return tc
	}

	if set {
		tc.want = testFlagValue
		tc.args = []string{"--" + testFlagName, testFlagValue}
	} else {
		tc.args = []string{"--" + boolFlagName}
	}

	return tc
}

func (tc *isSetTestcase) test(t *testing.T) {
	merged = false
	rawArgs = []string{}

	flag.CommandLine = flag.NewFlagSet("nodef:"+tc.name, flag.ContinueOnError)
	pflag.CommandLine = pflag.NewFlagSet("nodef:"+tc.name, pflag.ContinueOnError)

	var want string
	var wset bool

	if tc.fpkg != fpNone {
		rawArgs = tc.args

		want = tc.want
		dval := "default-value"

		if want == "" {
			want = dval
		} else {
			wset = true
		}

		switch tc.fpkg {
		case fpBase:
			flag.String(testFlagName, dval, "A test flag")
			flag.Bool(boolFlagName, false, "A bool flag")

		case fpSPF13:
			pflag.String(testFlagName, dval, "A test flag")
			pflag.Bool(boolFlagName, false, "A bool flag")
		}
	}

	if got, set, err := ValueIsSet(testFlagName); err != nil || got != want || set != wset {
		t.Errorf("ValueIsSet(%q) == (%q, %v, %v); Wanted(%q, %v, %v)", testFlagName, got, set, err, want, wset, nil)
	}
}

func TestValueIsSet(t *testing.T) {
	tests := []*isSetTestcase{
		newIsSetTestcase(fpNone, false),
		newIsSetTestcase(fpBase, false),
		newIsSetTestcase(fpSPF13, false),
		newIsSetTestcase(fpBase, true),
		newIsSetTestcase(fpSPF13, true),
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.test)
	}
}

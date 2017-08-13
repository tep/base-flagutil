package flagutil

import (
	"flag"
	"testing"

	"github.com/spf13/pflag"
)

type setUnlessTestcase struct {
	name string
	fpkg flagPackage
	fval string
	sval string
	tval string
	want bool
}

func (tc *setUnlessTestcase) test(t *testing.T) {
	var testvar string
	dval := "default-value"

	switch tc.fpkg {
	case fpBase:
		flag.StringVar(&testvar, testFlagName, dval, "A test flag")

	case fpSPF13:
		pflag.StringVar(&testvar, testFlagName, dval, "A test flag")
	}

	if tc.fval != "" {
		cmdlineArgs = []string{"--" + testFlagName, tc.fval}
	}

	if err := MergeAndParse(); err != nil {
		t.Fatalf("merging and parsing flags: %v", err)
	}

	if tc.sval != "" {
		if got, err := SetUnless(testFlagName, tc.sval); err != nil || got != tc.want {
			t.Errorf("SetUnless(%q, %q) == (%v, %v); Wanted (%v, %v)", testFlagName, tc.sval, got, err, tc.want, nil)
		}
	}

	if testvar != tc.tval {
		t.Errorf("testvar = %q; Wanted %q", testFlagName, tc.sval, testvar, tc.tval)
	}
}

func TestSetUnless(t *testing.T) {
	tests := []*setUnlessTestcase{
		{"noflags", fpNone, "", "", "", false},
		{"base-nocall", fpBase, "", "", "default-value", false},
		{"base-notset", fpBase, "", "override", "override", true},
		{"base-clset", fpBase, "cl-value", "override", "cl-value", false},
		{"spf13-nocall", fpSPF13, "", "", "default-value", false},
		{"spf13-notset", fpSPF13, "", "override", "override", true},
		{"spf13-clset", fpSPF13, "cl-value", "override", "cl-value", false},
	}

	for _, tc := range tests {
		resetForTest(tc.name)
		t.Run(tc.name, tc.test)
	}
}

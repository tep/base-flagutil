package flagutil

import (
	"flag"
	"testing"

	"github.com/spf13/pflag"
)

func TestMergeFlagSets(t *testing.T) {
	// TODO(tep): Flesh this out and make it do something.
}

type mergeFlagSetsTestcase struct {
	name   string
	from   *flag.FlagSet
	to     *pflag.FlagSet
	args   []string
	values map[string]string
}

func (tc *mergeFlagSetsTestcase) test(t *testing.T) {
	resetForTest(tc.name)
	flag.CommandLine = tc.from
	pflag.CommandLine = tc.to

	cmdlineArgs = tc.args

	if err := MergeAndParse(); err != nil {
		t.Fatal(err)
	}

	for k, v := range tc.values {
		fp := tc.to.Lookup(k)
		if fp == nil {
			t.Errorf("expected flag %q does not exist", k)
			continue
		}
		if fv := fp.Value.String(); fv != v {
			t.Errorf("flag value mismatch: Got %q; Wanted %q", fv, v)
		}
	}
}


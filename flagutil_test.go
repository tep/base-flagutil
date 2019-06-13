// Copyright © 2017 Tim Peoples <coders@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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

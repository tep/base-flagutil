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
	pFlagSet = pflag.CommandLine
}

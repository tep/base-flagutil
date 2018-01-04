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

// TODO(tep): Move this to "toolman.org/flags/flagsgroup"

package flagutil // import "toolman.org/base/flagutil"

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	FlagSetName     = os.Args[0]
	CommandLineArgs = os.Args[1:]
	GoFlagSet       = flag.CommandLine
	PFlagSet        = pflag.CommandLine
	debug           = false
	timenow         = time.Now
	pid             = os.Getpid()
)

func init() {
	debug, _ = strconv.ParseBool(os.Getenv("FLAGUTIL_DEBUG"))
	mergeGoFlagSet(true, PFlagSet, GoFlagSet)
}

func mergeGoFlagSet(hidden bool, pfs *pflag.FlagSet, gfs *flag.FlagSet) {
	pfs.AddGoFlagSet(gfs)
	if hidden {
		gfs.VisitAll(func(f *flag.Flag) {
			pfs.MarkHidden(f.Name)
		})
	}
}

//----------------------------------------------------------------------------

// FlagsGroup is a collection of FlagSets (from either the standard library or
// github.com/spf13/pflag) that may be merged and kept in sync as needed.
type FlagsGroup struct {
	fsm  fsMap
	base *pflag.FlagSet
}

// NewFlagsGroup creates a new FlagsGroup that includes the default FlagSet
// from both the standard flag library and github.com/spf13/pflag.
//
// Note that all flags defined by the standard flag library are marked as
// Hidden in this FlagsGroup's base FlagSet (but they should still function
// properly). See pflag.MarkHidden for more info.
func NewFlagsGroup() *FlagsGroup {
	fm := newSetterMap()

	fm.addGoFlagSet(GoFlagSet)
	fm.addFlagSet(PFlagSet)

	pfs := pflag.NewFlagSet(FlagSetName, pflag.ContinueOnError)
	pfs.AddFlagSet(PFlagSet)
	pfs.SortFlags = PFlagSet.SortFlags

	return &FlagsGroup{fm, pfs}
}

// SetPrimary sets the primary FlagSet to be p instead of the default created
// by NewFlagsGroup.
func (g *FlagsGroup) SetPrimary(p *pflag.FlagSet) {
	sf := p.SortFlags
	p.AddFlagSet(g.base)
	g.base = p
	g.base.SortFlags = sf
}

// AddFlagSet adds a new pflag.FlagSet to this FlagsGroup by first adding the
// new FlagSet to the FlagsGroup base and then adding the FlagsGroup base to
// the new FlagSet.
// See pflag.AddFlagSet for details about flag conflicts.
func (g *FlagsGroup) AddFlagSet(o *pflag.FlagSet) {
	g.base.AddFlagSet(o)
	g.fsm.addFlagSet(o)
	o.AddFlagSet(g.base)
}

// AddGoFlagSet adds a new flag.FlagSet (from the standard library) to this
// FlagsGroup.  Unlike the default FlagSet, these flags are not marked as
// Hidden in the FlagsGroup's base FlagSet.
func (g *FlagsGroup) AddGoFlagSet(o *flag.FlagSet) {
	mergeGoFlagSet(false, g.base, o)
}

// VisitAll exposes the method of the same name on the FlagsGroup base FlagSet.
func (g *FlagsGroup) VisitAll(vf func(*pflag.Flag)) {
	g.base.VisitAll(vf)
}

// Parse is a convenience wrapper around ParseArgs(os.Args[1:])
func (g *FlagsGroup) Parse() error {
	return g.ParseArgs(CommandLineArgs)
}

// ParseArgs calls Parse on the base FlagSet with the given args. It then
// merges the newly set flag values into each of the FlagsGroup's other
// FlagSets.
func (g *FlagsGroup) ParseArgs(args []string) error {
	g.base.Parse(args)

	if err := g.fsm.merge(g.base); err != nil {
		return err
	}

	GoFlagSet.Parse(g.base.Args())
	PFlagSet.Parse(g.base.Args())

	debugf("##### Flags Parsed and Merged")

	return nil
}

// ValueIsSet returns the value of the flag 'name' and a boolean indicating
// whether the value was modified by the parsed command line arguments.
func (g *FlagsGroup) ValueIsSet(name string) (string, bool) {
	f := g.base.Lookup(name)
	if f == nil {
		return "", false
	}

	return f.Value.String(), g.base.Changed(name)
}

// Set is used to set flag 'name' to 'val' in all included FlagSets. If the
// value of 'val' cannot be applied to this flag, an error is returned.
func (g *FlagsGroup) Set(name, val string) error {
	if err := g.base.Set(name, val); err != nil {
		return err
	}

	debugf("Setting %s=%q", name, val)
	if err := g.fsm.set(name, val); err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------

type flagSet struct {
	pfs  *pflag.FlagSet
	gfs  *flag.FlagSet
	name string
}

func (fs *flagSet) set(name, value string) error {
	if fs.pfs != nil {
		return fs.pfs.Set(name, value)
	}

	return fs.gfs.Set(name, value)
}

func (fs *flagSet) get(name string) string {
	if fs.pfs != nil {
		f := fs.pfs.Lookup(name)
		if f == nil {
			return ""
		}
		return f.Value.String()
	}

	f := fs.gfs.Lookup(name)
	if f == nil {
		return ""
	}
	return f.Value.String()
}

//----------------------------------------------------------------------------

type fsMap map[string][]*flagSet

func newSetterMap() fsMap {
	fm := fsMap(make(map[string][]*flagSet))
	return fm
}

func (fm *fsMap) addGoFlagSet(gfs *flag.FlagSet) {
	fs := &flagSet{gfs: gfs, name: identFlagSet(gfs)}
	gfs.VisitAll(func(f *flag.Flag) {
		(*fm)[f.Name] = append((*fm)[f.Name], fs)
	})
}

func (fm *fsMap) addFlagSet(pfs *pflag.FlagSet) {
	fs := &flagSet{pfs: pfs, name: identFlagSet(pfs)}
	pfs.VisitAll(func(f *pflag.Flag) {
		(*fm)[f.Name] = append((*fm)[f.Name], fs)
	})
}

func (fm *fsMap) merge(fs *pflag.FlagSet) error {
	vmap := make(map[string]string)
	fs.Visit(func(f *pflag.Flag) {
		vmap[f.Name] = f.Value.String()
	})

	for n, v := range vmap {
		debugf("Merging --%s=%q", n, v)
		if err := fm.set(n, v); err != nil {
			return err
		}
	}

	return nil
}

func (fm *fsMap) set(name, val string) error {
	for _, fs := range (*fm)[name] {
		if v := fs.get(name); v == val {
			debugf("    skipping --%s=%q for %s", name, val, fs.name)
			continue
		}

		debugf("    --%s=%q -> %s", name, val, fs.name)
		if err := fs.set(name, val); err != nil {
			return err
		}
	}

	return nil
}

//----------------------------------------------------------------------------

func debugf(msg string, args ...interface{}) {
	if !debug {
		return
	}

	var file string
	var line int
	var ok bool

	if _, file, line, ok = runtime.Caller(1); ok {
		if s := strings.LastIndex(file, "/"); s >= 0 {
			file = file[s+1:]
		}
	} else {
		file = "???"
		line = 1
	}

	fmt.Fprintf(os.Stderr, "D%s %7d %s:%d] %s\n",
		timenow().Format("0102 15:04:05.000000"),
		pid, file, line, fmt.Sprintf(msg, args...))
}

//----------------------------------------------------------------------------

func identFlagSet(i interface{}) string {
	var t string
	var v reflect.Value

	switch i.(type) {
	case *flag.FlagSet:
		t = "flag"
		v = reflect.ValueOf(i).Elem()

	case *pflag.FlagSet:
		t = "pflag"
		v = reflect.ValueOf(i).Elem()

	case flag.FlagSet:
		t = "flag"
		v = reflect.ValueOf(i)

	case pflag.FlagSet:
		t = "pflag"
		v = reflect.ValueOf(i)
	}

	if t == "" {
		return t
	}

	return t + ":" + v.FieldByName("name").String()
}

//----------------------------------------------------------------------------

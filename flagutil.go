package flagutil // import "toolman.org/base/flagutil"

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/pflag"
	"toolman.org/encoding/base56"
)

// TODO(tep): Write More Tests!!

var (
	merged      = make(map[string]bool)
	cmdlineArgs []string

	baseParsedError  = errors.New("cannot merge flags from package \"flag\": package already parsed")
	spf13ParsedError = errors.New("cannot merge flags into package \"github.com/spf13/pflag\": package already parsed")
)

func init() {
	// cmdlineArgs is used to override command line args in unit tests,
	// otherwise it's the real command-line arguments from os.Args
	cmdlineArgs = os.Args[1:]
}

func mergeKey(from *flag.FlagSet, to *pflag.FlagSet) string {
	b56 := func(i interface{}) string { return base56.Encode(uint64(reflect.ValueOf(i).Pointer())) }
	return fmt.Sprintf("%s:%s", b56(from), b56(to))
}

// Merged returns a boolean indicating whether the default FlagSets
// (i.e.  flag.CommandLine and pflag.CommandLine) have already been merged.
func Merged() bool {
	return FlagSetsMerged(flag.CommandLine, pflag.CommandLine)
}

// FlagSetsMerged returns a boolean indicating whether the given FlagSets
// have already been merged.
func FlagSetsMerged(from *flag.FlagSet, to *pflag.FlagSet) bool {
	return merged[mergeKey(from, to)]
}

// MergeFlags is a convenience wrapper around MergeFlagSets(flag.CommandLine, pflag.CommandLine).
func MergeFlags() error {
	return MergeFlagSets(flag.CommandLine, pflag.CommandLine)
}

// MergeFlagSets merges flags in the from flag.FlagSet into the to
// pflag.FlagSet.  If either FlagSet has previously been parsed, or if
// any flag names conflict between the two FlagSets, error is returned.
func MergeFlagSets(from *flag.FlagSet, to *pflag.FlagSet) error {
	mkey := mergeKey(from, to)
	if merged[mkey] {
		return nil
	}

	if from.Parsed() {
		return baseParsedError
	}

	if to.Parsed() {
		return spf13ParsedError
	}

	var nc []string

	from.VisitAll(func(f *flag.Flag) {
		if xf := to.Lookup(f.Name); xf != nil {
			nc = append(nc, f.Name)
			return
		}
		to.AddGoFlag(f)
	})

	if len(nc) > 0 {
		return newMergeConflictError(nc)
	}

	merged[mkey] = true
	return nil
}

func MergeAndParse() error {
	if err := MergeFlags(); err != nil {
		return err
	}

	if err := pflag.CommandLine.Parse(cmdlineArgs); err != nil {
		return err
	}

	return nil
}

// MergeConflictError is returned by MergeFlagSets or MergeFlags if the
// FlagSets being merged share flags of the same name.
type MergeConflictError struct {
	names []string
}

func newMergeConflictError(flags []string) *MergeConflictError {
	return &MergeConflictError{names: flags}
}

// IsMergeConflictError returns a boolean indicating whether the given error is
// a MergeConflictError error.
func IsMergeConflictError(err error) bool {
	_, ok := err.(*MergeConflictError)
	return ok
}

// Conflicts returns the list of common flag names that triggered the
// MergeConflictError.
func (e *MergeConflictError) Conflicts() []string {
	return e.names
}

func (e *MergeConflictError) Error() string {
	suf := ""
	if len(e.names) > 1 {
		suf = "s"
	}
	return fmt.Sprintf("name conflict%s merging flags: %v", suf, e.names)
}

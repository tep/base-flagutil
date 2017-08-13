package flagutil // import "toolman.org/base/flagutil"

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// TODO(tep): Write More Tests!!

var merged bool

func Merged() bool {
	return merged
}

func MergeFlags() error {
	return MergeFlagSets(flag.CommandLine, pflag.CommandLine)
}

func MergeFlagSets(from *flag.FlagSet, to *pflag.FlagSet) error {
	if merged {
		return nil
	}

	if from.Parsed() {
		return errors.New("cannot merge flags from package \"flag\": package already parsed")
	}

	if to.Parsed() {
		return errors.New("cannot merge flags into package \"github.com/spf13/pflag\": package already parsed")
	}

	var cnf []string

	from.VisitAll(func(f *flag.Flag) {
		if pf := to.Lookup(f.Name); pf != nil {
			cnf = append(cnf, f.Name)
			return
		}
		to.AddGoFlag(f)
	})

	if len(cnf) > 0 {
		sfx := ""
		if len(cnf) > 1 {
			sfx = "s"
		}
		return fmt.Errorf("name conflict%s merging flags: %v", sfx, cnf)
	}

	merged = true
	return nil
}

// Used by unit tests to override command line arguments
var cmdlineArgs []string

func init() {
	cmdlineArgs = os.Args[1:]
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

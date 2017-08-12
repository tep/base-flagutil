package flagutil // import "toolman.org/base/flagutil"

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// TODO(tep): Write Tests!!

var merged bool

func Merged() bool {
	return merged
}

func SetUnless(name, value string) (bool, error) {
	if _, ok, err := ValueIsSet(name); err != nil {
		return false, err
	} else if ok {
		return false, nil
	}

	if err := pflag.Set(name, value); err != nil {
		return false, err
	}

	return true, nil
}

func MergeFlags() error {
	if merged {
		return nil
	}

	if flag.Parsed() {
		return errors.New("cannot merge flags from package \"flag\": package already parsed")
	}

	if pflag.Parsed() {
		return errors.New("cannot merge flags into package \"github.com/spf13/pflag\": package already parsed")
	}

	var cnf []string

	flag.VisitAll(func(f *flag.Flag) {
		if pf := pflag.CommandLine.Lookup(f.Name); pf != nil {
			cnf = append(cnf, f.Name)
			return
		}
		pflag.CommandLine.AddGoFlag(f)
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

// Used by unit tests to override command line aarguments
var rawArgs []string

func init() {
	rawArgs = os.Args[1:]
}

func MergeAndParse() error {
	if err := MergeFlags(); err != nil {
		return err
	}

	if err := pflag.CommandLine.Parse(rawArgs); err != nil {
		return err
	}

	return nil
}

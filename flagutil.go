package flagutil // import "toolman.org/base/flagutil"

import (
	"errors"
	"flag"
	"fmt"

	"github.com/spf13/pflag"
)

// TODO(tep): Write Tests!!

var merged bool

func Merged() bool {
	return merged
}

func ValueIsSet(name string) (string, bool, error) {
	if err := MergeAndParse(); err != nil {
		return "", false, err
	}

	var isSet bool
	value := pflag.Lookup(name).Value.String()

	pflag.Visit(func(pf *pflag.Flag) {
		if pf.Name == name {
			isSet = true
			value = pf.Value.String()
		}
	})

	return value, isSet, nil
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

func MergeAndParse() error {
	if err := MergeFlags(); err != nil {
		return err
	}

	pflag.Parse()

	return nil
}

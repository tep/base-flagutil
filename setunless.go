package flagutil

import "github.com/spf13/pflag"

// TODO(tep): Write Unit Tests!!

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

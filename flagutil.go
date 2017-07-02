package flagutil  // import "toolman.org/base/flagutil"

import (
	"flag"
)

// TODO(tep): Write Tests!!

func ValueIsSet(name string) (string, bool) {
	if !flag.Parsed() {
		flag.Parse()
	}

	var isSet bool
	value := flag.Lookup(name).Value.String()

  flag.Visit(func(f *flag.Flag){
    if f.Name == name {
      isSet = true
			value = f.Value.String()
    }
  })

	return value, isSet
}

func SetUnless(name, value string) (bool, error) {
	if _, ok := ValueIsSet(name); ok {
		return false, nil
	}

  if err := flag.Set(name, value); err != nil {
    return false, err
  }

  return true, nil
}

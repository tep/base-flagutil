package flagutil

/*
import "github.com/spf13/pflag"

func ValueIsSet(name string) (string, bool, error) {
	if err := MergeAndParse(); err != nil {
		return "", false, err
	}

	var isSet bool
	fs := pflag.Lookup(name)
	if fs == nil {
		return "", false, nil
	}

	value := fs.DefValue

	pflag.Visit(func(pf *pflag.Flag) {
		if pf.Name == name {
			isSet = true
			value = pf.Value.String()
		}
	})

	return value, isSet, nil
}
*/

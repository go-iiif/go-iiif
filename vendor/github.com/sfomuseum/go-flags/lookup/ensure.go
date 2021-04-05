package lookup

import (
	"flag"
)

func EnsureStringVars(fs *flag.FlagSet, vars ...string) error {

	for _, v := range vars {

		_, err := StringVar(fs, v)

		if err != nil {
			return err
		}
	}

	return nil
}

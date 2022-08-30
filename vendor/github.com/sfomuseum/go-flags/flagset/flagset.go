// package flagset provides methods for working with `flag.FlagSet` instances.
package flagset

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Parse command line arguments with a flag.FlagSet instance.
func Parse(fs *flag.FlagSet) {

	args := os.Args[1:]

	if len(args) > 0 && args[0] == "-h" {
		fs.Usage()
		os.Exit(0)
	}

	fs.Parse(args)
}

// Assign values to a flag.FlagSet instance from matching environment variables.
func SetFlagsFromEnvVars(fs *flag.FlagSet, prefix string) error {
	return SetFlagsFromEnvVarsWithFeedback(fs, prefix, false)
}

// Assign values to a flag.FlagSet instance from matching environment variables, optionally logging progress and other feedback.
func SetFlagsFromEnvVarsWithFeedback(fs *flag.FlagSet, prefix string, feedback bool) error {

	fs.VisitAll(func(fl *flag.Flag) {

		name := fl.Name
		env := FlagNameToEnvVar(prefix, name)

		val, ok := os.LookupEnv(env)

		if ok {

			if feedback {
				log.Printf("set -%s flag from %s environment variable\n", name, env)
			}

			fs.Set(name, val)
		}
	})

	return nil
}

// Create a new flag.FlagSet instance.
func NewFlagSet(name string) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ExitOnError)

	fs.Usage = func() {
		fs.PrintDefaults()
	}

	return fs
}

// FlagNameToEnvVar formats 'name' and 'prefix' in to an environment variable name, used to lookup
// a value.
func FlagNameToEnvVar(prefix string, name string) string {

	prefix = normalizeEnvVar(prefix)
	name = normalizeEnvVar(name)

	return fmt.Sprintf("%s_%s", prefix, name)

}

// normalizeEnvVar normalizes a flag name in to its corresponding environment variable name.
func normalizeEnvVar(raw string) string {

	new := raw

	new = strings.ToUpper(new)
	new = strings.Replace(new, "-", "_", -1)

	return new
}

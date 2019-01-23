package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func SetFlagsFromEnvVars(prefix string) error {

	prefix = normalize(prefix)

	flag.VisitAll(func(fl *flag.Flag) {

		name := fl.Name

		env_name := fmt.Sprintf("%s_%s", prefix, normalize(name))

		v, ok := os.LookupEnv(env_name)

		if ok {

			flag.Set(name, v)
			log.Printf("Reset %s flag from %s environment variable\n", name, env_name)
		}
	})

	return nil

}

func normalize(s string) string {

	s = strings.Replace(s, "-", "_", -1)
	s = strings.ToUpper(s)

	return s
}

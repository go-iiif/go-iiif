package flags

import (
	"flag"
	"github.com/go-ini/ini"
	"log"
)

func SetFlagsFromConfig(path string, section string) error {

	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys: true,
	}, path)

	if err != nil {
		return err
	}

	sect, err := cfg.GetSection(section)

	if err != nil {
		return err
	}

	flag.VisitAll(func(fl *flag.Flag) {

		name := fl.Name

		if sect.HasKey(name) {

			k := sect.Key(name)
			v := k.Value()

			flag.Set(name, v)

			log.Printf("Reset %s flag from config file\n", name)
		}
	})

	return nil
}

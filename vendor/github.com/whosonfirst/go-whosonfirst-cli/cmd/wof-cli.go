package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
)

func main() {

	var dsn flags.MultiDSNString
	flag.Var(&dsn, "dsn", "...")

	flag.Parse()

	log.Println(dsn.String())
}

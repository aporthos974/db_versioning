package main

import (
	"db_versioning/initialisation"
	"db_versioning/migration"
	"db_versioning/version"
	"flag"
	"fmt"
	"os"
)

func main() {
	var initialize, upgrade, displayVersion = initArgsAndFlags()

	checkParameters()

	flag.Visit(func(f *flag.Flag) {
		schema := flag.Arg(0)
		if f.Name == "i" && initialize {
			initialisation.Initialize(schema)
		} else if f.Name == "v" && displayVersion {
			fmt.Printf("current version : %s \n", version.GetCurrentVersion())
		} else if f.Name == "u" && upgrade {
			migration.Migrate(schema)
		}
	})
}

func initArgsAndFlags() (bool, bool, bool) {
	var initialize, upgrade, displayVersion bool
	flag.BoolVar(&initialize, "i", false, "Initialize versioning system for database schema")
	flag.BoolVar(&upgrade, "u", false, "Upgrade database schema")
	flag.BoolVar(&displayVersion, "v", false, "Display database schema version")
	return initialize, upgrade, displayVersion
}

func checkParameters() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Missing schema argument \n")
		fmt.Printf("Usage of %s [option] <schema> \n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	if flag.NFlag() == 0 {
		fmt.Printf("Missing flag \n")
		fmt.Printf("Usage of %s [option] <schema> \n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
}

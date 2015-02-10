package main

import (
	"db_versioning/initialisation"
	"db_versioning/migration"
	"db_versioning/reinit"
	"db_versioning/version"
	"flag"
	"fmt"
	"os"
)

func main() {
	var initialize, reinitialize, upgrade, displayVersion = initArgsAndFlags()

	checkParameters()

	flag.Visit(func(f *flag.Flag) {
		schema := flag.Arg(0)
		fmt.Printf("_______________________________________________ \n")
		if f.Name == "i" && *initialize {
			fmt.Printf("\nInitialize database schema version... \n")
			initialisation.Initialize(schema)
		} else if f.Name == "v" && *displayVersion {
			fmt.Printf("\nGet current version... \n")
			version.DisplayCurrentVersion(schema)
		} else if f.Name == "u" && *upgrade {
			fmt.Printf("\nUpdate database schema... \n")
			migration.Migrate(schema)
			version.DisplayCurrentVersion(schema)
		} else if f.Name == "I" && *reinitialize {
			fmt.Printf("\nRe-initialize database schema... \n")
			reinit.Reinitialize(schema)
		}
	})
}

func initArgsAndFlags() (*bool, *bool, *bool, *bool) {
	var initialize, reinitialize, upgrade, displayVersion bool
	var environment string
	flag.BoolVar(&initialize, "i", false, "Initialize versioning system for database schema")
	flag.BoolVar(&reinitialize, "I", false, "Delete and create database schema, initialize versioning system and upgrade")
	flag.BoolVar(&upgrade, "u", false, "Upgrade database schema")
	flag.BoolVar(&displayVersion, "v", false, "Display database schema version")
	flag.StringVar(&environment, "host", "localhost", "Database environment (not implemented)")
	return &initialize, &reinitialize, &upgrade, &displayVersion
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

package main

import (
	"db_versioning/db"
	"db_versioning/initialisation"
	"db_versioning/migration"
	"db_versioning/reinit"
	"db_versioning/version"
	"flag"
	"fmt"
	"os"
)

type FlagValues struct {
	Initialize, Reinitialize, Upgrade, DisplayVersion bool
	Environment                                       string
}

func main() {
	var flagValues = initArgsAndFlags()

	checkParameters()

	db.Host = flagValues.Environment

	flag.Visit(func(f *flag.Flag) {
		schema := flag.Arg(0)
		if f.Name != "host" {
			fmt.Printf("_______________________________________________ \n")
		}
		if f.Name == "i" && flagValues.Initialize {
			fmt.Printf("\nInitialize database schema version... \n")
			initialisation.Initialize(schema)
		} else if f.Name == "v" && flagValues.DisplayVersion {
			fmt.Printf("\nGet current version... \n")
			version.DisplayCurrentVersion(schema)
		} else if f.Name == "u" && flagValues.Upgrade {
			fmt.Printf("\nUpdate database schema... \n")
			migration.Migrate(schema)
			version.DisplayCurrentVersion(schema)
		} else if f.Name == "I" && flagValues.Reinitialize {
			fmt.Printf("\nRe-initialize database schema... \n")
			reinit.Reinitialize(schema)
		}
	})
}

func initArgsAndFlags() *FlagValues {
	var flagValues FlagValues
	flag.BoolVar(&flagValues.Initialize, "i", false, "Initialize versioning system for database schema")
	flag.BoolVar(&flagValues.Reinitialize, "I", false, "Delete and create database schema, initialize versioning system and upgrade")
	flag.BoolVar(&flagValues.Upgrade, "u", false, "Upgrade database schema")
	flag.BoolVar(&flagValues.DisplayVersion, "v", false, "Display database schema version")
	flag.StringVar(&flagValues.Environment, "host", "localhost", "Database environment (not implemented)")
	return &flagValues
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

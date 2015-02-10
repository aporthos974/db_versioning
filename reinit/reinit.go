package reinit

import (
	"db_versioning/initialisation"
	"db_versioning/migration"
)

func Reinitialize(schema string) {
	migration.ExecuteQueries("", "drop database "+schema, "create database "+schema)
	initialisation.Initialize(schema)
	migration.Migrate(schema)
}

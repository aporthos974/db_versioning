package migration

import (
	"db_versioning/db"
	"db_versioning/version"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanApplyScript(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("db_versioning_test")

	assert.Equal(test, "1.0.10", version.GetCurrentVersion("db_versioning_test"))
}

func TestCanApplySeveralScripts(test *testing.T) {
	db.InitDatabase("0.0.1")

	Migrate("db_versioning_test")

	versions := db.GetVersions()
	assert.Equal(test, "0.0.1", versions[0].Version)
	assert.Equal(test, "1.0.0", versions[1].Version)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "1.0.2", versions[3].Version)
	assert.Equal(test, "1.0.10", versions[4].Version)
	assert.Equal(test, "1.0.10", version.GetCurrentVersion("db_versioning_test"))
}

func TestCanApplySeveralScriptsFromVersion(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("db_versioning_test")

	assert.Equal(test, "1.0.10", version.GetCurrentVersion("db_versioning_test"))
}

func TestCanApplySeveralScriptsInTheSameVersion(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("db_versioning_test")

	versions := db.GetVersions()
	assert.Len(test, versions, 4)
	assert.Equal(test, "1.0.0", versions[0].Version)
	assert.Equal(test, "1.0.1", versions[1].Version)
	assert.Equal(test, "db_versioning_test/1.0.1/first.sql;db_versioning_test/1.0.1/second.sql", versions[1].Script)
	assert.Equal(test, "1.0.2", versions[2].Version)
	assert.Equal(test, "1.0.10", versions[3].Version)
	assert.Equal(test, "1.0.10", version.GetCurrentVersion("db_versioning_test"))
}

func TestCanKnownSchemaIsAlreadyUpToDate(test *testing.T) {
	db.InitDatabase("1.0.1")

	Migrate("db_versioning_test")

	versions := db.GetVersions()
	assert.Len(test, versions, 3)
	assert.Equal(test, "1.0.1", versions[0].Version)
	assert.Equal(test, "1.0.2", versions[1].Version)
	assert.Equal(test, "1.0.10", versions[2].Version)
	assert.Equal(test, "1.0.10", version.GetCurrentVersion("db_versioning_test"))
}

func TestCanApplyScriptsInGoodOrder(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("db_versioning_test")

	versions := db.GetVersions()
	assert.Len(test, versions, 4)
	assert.Equal(test, "1.0.0", versions[0].Version)
	assert.Equal(test, "1.0.1", versions[1].Version)
	assert.Equal(test, "1.0.2", versions[2].Version)
	assert.Equal(test, "1.0.10", versions[3].Version)
	assert.Equal(test, "1.0.10", version.GetCurrentVersion("db_versioning_test"))
}

func TestMigrationDoesntLaunchWhenMigrationWasAlreadyFailed(test *testing.T) {
	db.InitDatabaseVersion("1.0.0", "failed")
	ExecuteQuery("insert into db_version (script, version, state) values ('test.sql', '0.0.0', 'ok')", "db_versioning_test")

	Migrate("db_versioning_test")

	versions := db.GetVersions()
	assert.Len(test, versions, 2)
	assert.Equal(test, "0.0.0", version.GetCurrentVersion("db_versioning_test"))
}

package migration

import (
	"db_versioning/db"
	"db_versioning/version"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanApplyScript(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("../db_versioning_test_ok")

	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScripts(test *testing.T) {
	db.InitDatabase("0.0.0")

	Migrate("../db_versioning_test_ok")

	versions := db.GetVersions()
	assert.Equal(test, "0.0.0", versions[0].Version)
	assert.Equal(test, "1.0.0", versions[1].Version)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScriptsFromVersion(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("../db_versioning_test_ok")

	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScriptsInTheSameVersion(test *testing.T) {
	db.InitDatabase("1.0.0")

	Migrate("../db_versioning_test_ok")

	versions := db.GetVersions()
	assert.Equal(test, 3, len(versions))
	assert.Equal(test, "1.0.0", versions[0].Version)
	assert.Equal(test, "1.0.1", versions[1].Version)
	assert.Equal(test, "../db_versioning_test_ok/1.0.1/first.sql", versions[1].Script)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "../db_versioning_test_ok/1.0.1/second.sql", versions[2].Script)
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanKnownScriptFailed(test *testing.T) {
	db.InitDatabase("0.0.0")

	assert.Panics(test, func() { Migrate("../db_versioning_test_failed") }, "Calling Compare() should panic")

	versions := db.GetVersions()
	assert.Equal(test, 3, len(versions))
	assert.Equal(test, "0.0.0", versions[0].Version)
	assert.Equal(test, "1.0.0", versions[1].Version)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "../db_versioning_test_failed/1.0.1/error.sql", versions[2].Script)
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

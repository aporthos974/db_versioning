package migration

import (
	"db_versioning/version"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Version struct {
	Version, Script string
}

func TestCanApplyScript(test *testing.T) {
	initDatabase("1.0.0")

	Migrate("../db_versioning_test_ok")

	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScripts(test *testing.T) {
	initDatabase("0.0.0")

	Migrate("../db_versioning_test_ok")

	versions := getVersions()
	assert.Equal(test, "0.0.0", versions[0].Version)
	assert.Equal(test, "1.0.0", versions[1].Version)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScriptsFromVersion(test *testing.T) {
	initDatabase("1.0.0")

	Migrate("../db_versioning_test_ok")

	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScriptsInTheSameVersion(test *testing.T) {
	initDatabase("1.0.0")

	Migrate("../db_versioning_test_ok")

	versions := getVersions()
	assert.Equal(test, 3, len(versions))
	assert.Equal(test, "1.0.0", versions[0].Version)
	assert.Equal(test, "1.0.1", versions[1].Version)
	assert.Equal(test, "../db_versioning_test_ok/1.0.1/first.sql", versions[1].Script)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "../db_versioning_test_ok/1.0.1/second.sql", versions[2].Script)
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanKnownScriptFailed(test *testing.T) {
	initDatabase("0.0.0")

	assert.Panics(test, func() { Migrate("../db_versioning_test_failed") }, "Calling Compare() should panic")

	versions := getVersions()
	assert.Equal(test, 3, len(versions))
	assert.Equal(test, "0.0.0", versions[0].Version)
	assert.Equal(test, "1.0.0", versions[1].Version)
	assert.Equal(test, "1.0.1", versions[2].Version)
	assert.Equal(test, "../db_versioning_test_failed/1.0.1/error.sql", versions[2].Script)
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func getVersions() []Version {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	rows, _, _ := db.Query("select version, script from db_version order by id")
	db.Close()
	var versions []Version
	for _, row := range rows {
		versions = append(versions, Version{Version: row.Str(0), Script: row.Str(1)})
	}
	return versions
}

func initDatabase(targetVersion string) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	dropAllTables(db)
	db.Query("truncate db_version")
	db.Query("insert into db_version (script, version, state) values ('test.sql', '%s', 'ok')", targetVersion)
	db.Close()
}

func dropAllTables(db mysql.Conn) {
	rows, _, _ := db.Query("show tables")
	var tables []string
	for _, row := range rows {
		table := row.Str(0)
		if table != "db_version" {
			tables = append(tables, table)
		}
	}
	concatenateTables := strings.Join(tables, ", ")
	db.Query("drop table " + concatenateTables)
}

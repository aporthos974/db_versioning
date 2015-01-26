package migration

import (
	"db_versioning/version"
	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"strings"
	"testing"
)

func TestCanApplyScript(test *testing.T) {
	initDatabase("1.0.0")

	Migrate("../db_versioning_test/")

	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func TestCanApplySeveralScripts(test *testing.T) {
	initDatabase("0.0.0")

	Migrate("../db_versioning_test/")

	versions := getVersions()
	assert.Equal(test, "0.0.0", versions[0])
	assert.Equal(test, "1.0.0", versions[1])
	assert.Equal(test, "1.0.1", versions[2])
	assert.Equal(test, "1.0.1", version.GetCurrentVersion())
}

func getVersions() []string {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	rows, _, _ := db.Query("select version from db_version order by id")
	db.Close()
	var versions []string
	for _, row := range rows {
		versions = append(versions, row.Str(0))
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

package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func TestCanGetCurrentDBVersion(test *testing.T) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	dropAllTables(db)
	db.Query("create table db_version (id INTEGER PRIMARY KEY AUTO_INCREMENT , script VARCHAR(255), version VARCHAR(255), state VARCHAR(255))")
	db.Query("insert into db_version (script, version, state) values ('test.sql', '1.0.0', 'ok')")
	db.Query("insert into db_version (script, version, state) values ('test.sql', '1.0.1', 'ok')")

	version := GetCurrentVersion()

	assert.Equal(test, "1.0.1", version)
}

func TestCanKnownFirstVersionIsGreaterToSecondVersion(test *testing.T) {

	compare := Compare("1.1.1", "1.0.0")

	assert.Equal(test, 1, compare)
}

func TestCanKnownFirstVersionIsLowerToSecondVersion(test *testing.T) {

	compare := Compare("1.0.1", "1.1.0")

	assert.Equal(test, -1, compare)
}

func TestCanKnownFirstVersionIsEqualToSecondVersion(test *testing.T) {

	compare := Compare("1.0.1", "1.0.1")

	assert.Equal(test, 0, compare)
}

func TestCanKnownFirstVersionFormatIsIncompatibleWithSecondVersionFormat(test *testing.T) {
	compareFunction := func() { Compare("1.0.1.1", "1.0.1") }

	assert.Panics(test, compareFunction, "Calling Compare() should panic")
}

func TestCanKnownFirstVersionFormatIsNotSupported(test *testing.T) {
	compareFunction := func() { Compare("1.0.1a", "1.0.1") }

	assert.Panics(test, compareFunction, "Calling Compare() should panic")
}

func TestCanKnownSecondVersionFormatIsNotSupported(test *testing.T) {
	compareFunction := func() { Compare("1.0.1", "1a.0.1") }

	assert.Panics(test, compareFunction, "Calling Compare() should panic")
}

func dropAllTables(db mysql.Conn) {
	rows, _, _ := db.Query("show tables")
	var tables []string
	for _, row := range rows {
		tables = append(tables, row.Str(0))
	}
	concatenateTables := strings.Join(tables, ", ")
	db.Query("drop table " + concatenateTables)
}

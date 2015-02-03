package version

import (
	"db_versioning/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func TestCanGetCurrentDBVersion(test *testing.T) {
	db.InitDatabase("0.0.0")
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	db.Query("insert into db_version (script, version, state) values ('test.sql', '1.0.0', 'ok')")
	db.Query("insert into db_version (script, version, state) values ('test.sql', '1.0.1', 'ok')")
	db.Close()

	version := GetCurrentVersion("db_versioning_test")

	assert.Equal(test, "1.0.1", version)
}

func TestCanSortVersions(test *testing.T) {
	versions := []string{"10.0.0", "1.0.10", "1.0.1", "1.0.2", "1.20.0", "1.2.1"}

	sortedVersions := Sort(versions)

	expected := []string{"1.0.1", "1.0.2", "1.0.10", "1.2.1", "1.20.0", "10.0.0"}
	assert.Equal(test, expected, sortedVersions)
}

func TestCanKnownFirstVersionIsGreaterToSecondVersion(test *testing.T) {

	compare := Compare("10.10.10", "1.0.0")

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

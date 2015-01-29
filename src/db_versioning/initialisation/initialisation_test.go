package initialisation

import (
	db_utils "db_versioning/db"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Version struct {
	Version, Script string
}

func TestCanInitWhenTableDBVersionDoesntExist(test *testing.T) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	dropAllTables(db)
	db.Query("drop table db_version")
	db.Close()

	Initialize()

	versions := db_utils.GetVersions()
	assert.Equal(test, 1, len(versions))
	assert.Equal(test, "0.0.0", versions[0].Version)
	assert.Equal(test, "initialisation", versions[0].Script)
}

func TestDoesntInitWhenTableDBVersionExists(test *testing.T) {
	db_utils.InitDatabase("0.0.0")

	Initialize()

	versions := db_utils.GetVersions()
	assert.Equal(test, 1, len(versions))
	assert.Equal(test, "0.0.0", versions[0].Version)
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

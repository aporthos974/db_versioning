package reinit

import (
	dbinit "db_versioning/db"
	"db_versioning/version"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func TestReinitDB(test *testing.T) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "")
	db.Connect()
	db.Query("create database if not exists db_versioning_test")
	db.Query("create table test_table (test_id VARCHAR(255), test_label VARCHAR(255))")
	db.Close()
	dbinit.InitDatabase("0.0.0")

	Reinitialize("db_versioning_test")

	version := version.GetCurrentVersion("db_versioning_test")
	assert.Equal(test, "1.0.10", version)
}

package version

import (
	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"testing"
)

func TestCanGetCurrentDBVersion(test *testing.T) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()
	db.Query("truncate db_version")
	db.Query("insert into db_version (script, version, state) values ('test.sql', '1.0.0', 'ok')")
	db.Query("insert into db_version (script, version, state) values ('test.sql', '1.0.1', 'ok')")

	version := GetCurrentVersion()

	assert.Equal(test, "1.0.1", version)
}

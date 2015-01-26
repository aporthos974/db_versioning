package version

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"log"
)

func GetCurrentVersion() string {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	err := db.Connect()
	if err != nil {
		log.Fatalf("Connection failed : %s", err.Error())
	}
	versionRow, _, err := db.QueryFirst("select version from db_version order by id desc limit 1")
	if err != nil {
		log.Fatalf("Query failed : %s", err.Error())
	}
	return versionRow.Str(0)
}

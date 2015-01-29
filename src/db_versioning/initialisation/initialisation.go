package initialisation

import (
	"log"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func Initialize() {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()

	row, _, err := db.QueryFirst("show tables like 'db_version'")
	if err != nil {
		log.Panicf("Error while fetching if table db_version exists")
	}
	if row == nil {
		db.Query("create table db_version (id INTEGER PRIMARY KEY AUTO_INCREMENT , script VARCHAR(255), version VARCHAR(255), state VARCHAR(255))")
		db.Query("insert into db_version (script, version, state) values ('initialisation', '0.0.0', 'ok')")
	}
	db.Close()
}

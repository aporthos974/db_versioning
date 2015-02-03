package log

import (
	"db_versioning/db"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	"os"
)

func FailAndLogInDatabase(db mysql.Conn, script db.Script, errorMessage string, err error) {
	UpgradeDBVersion(script.Version, script.Path, fmt.Sprintf("failed : %s", err.Error()), db)
	Fail(errorMessage, err.Error())
}

func Fail(errorMessage string, args ...interface{}) {
	fmt.Printf(errorMessage, args...)
	os.Exit(1)
}

func UpgradeDBVersion(toVersion, scriptName string, state string, db mysql.Conn) {
	statement, err := db.Prepare("insert into db_version (script, version, state) values (?, ?, ?)")
	if err != nil {
		Fail("Error while preparing the update db_version : %s \n", err.Error())
	}

	_, err = statement.Run(scriptName, toVersion, state)
	if err != nil {
		Fail("Error while updating db_version : %s \n", err.Error())
	}
}

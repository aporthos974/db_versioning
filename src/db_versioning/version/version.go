package version

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"log"
	"regexp"
	"strings"
)

func Compare(firstVersion string, secondVersion string) int {
	if isFormatValid(firstVersion) || isFormatValid(secondVersion) {
		panic(fmt.Sprintf("Error incompatible version format : %s / %s", firstVersion, secondVersion))
	}

	firstSplittedVersion := strings.Split(firstVersion, ".")
	secondSplittedVersion := strings.Split(secondVersion, ".")
	for i := 0; i < len(firstSplittedVersion); i++ {
		if firstSplittedVersion[i] > secondSplittedVersion[i] {
			return 1
		} else if firstSplittedVersion[i] < secondSplittedVersion[i] {
			return -1
		}
	}
	return 0
}

func isFormatValid(version string) bool {
	valid, err := regexp.MatchString("^\\d.\\d.\\d$", version)
	return !valid || err != nil
}

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

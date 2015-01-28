package migration

import (
	"bytes"
	"db_versioning/version"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Script struct {
	Path, Version string
	Queries       []string
}

func Migrate(schemaPath string) {
	folders, _ := ioutil.ReadDir(schemaPath)
	var scripts []Script
	for _, folder := range folders {
		if folder.IsDir() && version.Compare(folder.Name(), version.GetCurrentVersion()) == 1 {
			files, _ := ioutil.ReadDir(computePath(schemaPath, folder.Name()))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".sql") {
					queries := fetchQueries(computePath(schemaPath, folder.Name(), file.Name()))
					scripts = append(scripts, Script{Path: computePath(schemaPath, folder.Name(), file.Name()), Version: folder.Name(), Queries: queries})
				}
			}
		}
	}

	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	db.Connect()

	executeScripts(scripts, db)
	db.Close()
}

func computePath(basePath string, elementsPath ...string) string {
	var path bytes.Buffer
	path.WriteString(basePath)
	for _, element := range elementsPath {
		path.WriteString("/")
		path.WriteString(element)
	}
	return path.String()
}

func fetchQueries(scriptPath string) []string {
	content, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("Error when openning file : %s", err.Error())
	}
	return strings.Split(string(content), ";")
}

func executeScripts(scripts []Script, db mysql.Conn) {
	for _, script := range scripts {
		for _, query := range script.Queries {
			if !isEmptyString(query) {
				_, _, err := db.Query(query)
				if err != nil {
					log.Fatalf("Error when executing script : %s", err.Error())
				}
			}
		}
		upgradeDBVersion(script.Version, script.Path, db)
	}
}

func isEmptyString(value string) bool {
	return strings.TrimSpace(value) == ""
}

func upgradeDBVersion(toVersion, scriptName string, db mysql.Conn) {
	statement, err := db.Prepare("insert into db_version (script, version, state) values (?, ?, 'ok')")
	if err != nil {
		log.Fatalf("Error when preparing the update db_version : %s", err.Error())
	}

	_, err = statement.Run(scriptName, toVersion)
	if err != nil {
		log.Fatalf("Error when updating db_version : %s", err.Error())
	}
}

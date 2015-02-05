package migration

import (
	"bytes"
	. "db_versioning/db"
	. "db_versioning/log"
	"db_versioning/version"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

func Migrate(schema string) {
	if rows, _ := executeQuery("select count(*) from db_version where state <> 'ok'", schema); rows[0].Int(0) > 0 {
		fmt.Printf("There is a script in error. Fix script, delete version in error in table 'db_version' and relaunch db_versioning command. \n")
		return
	}
	scripts := fetchMigrationScripts(schema)
	if len(scripts) == 0 {
		fmt.Printf("Schema is already up-to-date \n")
		return
	}
	executeScripts(scripts, schema)
	fmt.Printf("Database schema '%s' updated \n", schema)
}

func executeQuery(query string, schema string) ([]mysql.Row, mysql.Result) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", schema)
	err := db.Connect()
	if err != nil {
		Fail("Connection failed : %s \n", err.Error())
	}
	rows, result, err := db.Query(query)
	db.Close()
	if err != nil {
		Fail("Query execution failed : %s \n", err.Error())
	}
	return rows, result
}

func fetchMigrationScripts(schema string) []Script {
	folders, _ := ioutil.ReadDir(schema)
	folders = sortDescFolders(folders)
	var scripts []Script
	currentVersion := version.GetCurrentVersion(schema)
	for _, folder := range folders {
		if isEligibleFolder(folder, currentVersion) {
			files, _ := ioutil.ReadDir(computePath(schema, folder.Name()))
			var queries []Query
			var scriptPaths []string
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".sql") {
					queries = append(queries, fetchQueries(computePath(schema, folder.Name(), file.Name()))...)
					scriptPaths = append(scriptPaths, computePath(schema, folder.Name(), file.Name()))
				}
			}
			scripts = append(scripts, createScript(scriptPaths, folder, queries))
		}
	}
	return sortAscScripts(scripts)
}

func sortDescFolders(folders []os.FileInfo) []os.FileInfo {
	sort.Sort(sort.Reverse(FolderSort(folders)))
	return folders
}

func sortAscScripts(scripts []Script) []Script {
	sort.Sort(ScriptSort(scripts))
	return scripts
}

func createScript(scriptPaths []string, folder os.FileInfo, queries []Query) Script {
	return Script{Paths: strings.Join(scriptPaths, ";"), Version: folder.Name(), Queries: queries}
}

func isEligibleFolder(folder os.FileInfo, currentVersion string) bool {
	return folder.IsDir() && version.Compare(folder.Name(), currentVersion) == 1
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

func fetchQueries(scriptPath string) []Query {
	content, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		Fail("Error when openning file : %s \n", err.Error())
	}
	var queries []Query
	for _, query := range strings.Split(string(content), ";") {
		queries = append(queries, Query(query))
	}
	return queries
}

func executeScripts(scripts []Script, schema string) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", schema)
	err := db.Connect()
	if err != nil {
		Fail("Connection failed : %s \n", err.Error())
	}
	for _, script := range scripts {
		transaction, err := db.Begin()
		if err != nil {
			FailAndLogInDatabase(db, script, "Error while opening transaction : %s \n", err)
		}
		for _, query := range script.Queries {
			if !query.IsEmpty() {
				_, _, err := transaction.Query(query.GetContent())
				if err != nil {
					transaction.Rollback()
					FailAndLogInDatabase(db, script, "Error while executing script : %s \n", err)
				}
			}
		}
		transaction.Commit()
		fmt.Printf("executed : %s \n", script.Paths)
		UpgradeDBVersion(script.Version, script.Paths, "ok", db)
	}
	db.Close()
}

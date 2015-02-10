package migration

import (
	. "db_versioning/db"
	. "db_versioning/log"
	"db_versioning/version"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

var EXECUTABLE_PATH = getExecutablePath()

func Migrate(schema string) {
	if rows, _ := ExecuteQuery("select count(*) from db_version where state <> 'ok'", schema); rows[0].Int(0) > 0 {
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

func ExecuteQueries(schema string, queries ...string) (results []mysql.Result) {
	db := openDBConnection("127.0.0.1:3306", "test", "test", schema)
	for _, query := range queries {
		_, result, err := db.Query(query)
		if err != nil {
			Fail("Query execution failed : %s \n", err.Error())
		}
		results = append(results, result)
	}
	db.Close()
	return results
}

func openDBConnection(host string, login string, passwd string, schema string) mysql.Conn {
	db := mysql.New("tcp", "", host, login, passwd, schema)
	err := db.Connect()
	if err != nil {
		Fail("Connection failed : %s \n", err.Error())
	}
	return db
}

func ExecuteQuery(query string, schema string) ([]mysql.Row, mysql.Result) {
	db := openDBConnection("127.0.0.1:3306", "test", "test", schema)
	rows, result, err := db.Query(query)
	db.Close()
	if err != nil {
		Fail("Query execution failed : %s \n", err.Error())
	}
	return rows, result
}

func fetchMigrationScripts(schema string) []Script {
	var scripts []Script
	currentVersion := version.GetCurrentVersion(schema)
	versionsDir := filepath.Join(EXECUTABLE_PATH, schema)
	for _, folder := range readVersionFolders(versionsDir) {
		versionDir := filepath.Join(versionsDir, folder.Name())
		if folder.IsDir() {
			if !isEligibleFolder(folder, currentVersion) {
				break
			}
			var queries []Query
			var scriptPaths []string
			for _, file := range readFilesInFolder(versionDir) {
				if isSQLType(file) {
					sqlScriptPath := filepath.Join(versionDir, file.Name())
					queries = append(queries, fetchQueries(sqlScriptPath)...)
					scriptPaths = append(scriptPaths, filepath.Join(schema, folder.Name(), file.Name()))
				}
			}
			scripts = append(scripts, createScript(scriptPaths, folder, queries))
		}
	}
	return sortAscScripts(scripts)
}

func isSQLType(file os.FileInfo) bool {
	return strings.HasSuffix(file.Name(), ".sql")
}

func readFilesInFolder(versionDir string) []os.FileInfo {
	files, err := ioutil.ReadDir(versionDir)
	if err != nil {
		Fail("Error while reading each files in folder : %s", err.Error())
	}
	return files
}

func readVersionFolders(versionsDir string) []os.FileInfo {
	folders, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		Fail("Error while reading version folder : %s \n", err.Error())
	}
	return sortDescFolders(folders)
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
	return version.Compare(folder.Name(), currentVersion) == 1
}

func getExecutablePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), "..")
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
	db := openDBConnection("127.0.0.1:3306", "test", "test", schema)
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

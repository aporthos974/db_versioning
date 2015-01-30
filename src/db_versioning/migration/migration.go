package migration

import (
	"bytes"
	"db_versioning/version"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Script struct {
	Path, Version string
	Queries       []Query
}

type Query string

func (query Query) isEmpty() bool {
	return strings.TrimSpace(query.GetContent()) == ""
}

func (query Query) GetContent() string {
	return fmt.Sprint(query)
}

func Migrate(schema string) {
	scripts := fetchMigrationScripts("../" + schema)
	if len(scripts) == 0 {
		fmt.Printf("Schema is already up-to-dateUS")
		return
	}
	executeScripts(scripts)
	fmt.Printf("Database schema '%s' updated \n", schema)
}

func fetchMigrationScripts(schemaPath string) []Script {
	folders, _ := ioutil.ReadDir(schemaPath)
	var scripts []Script
	for _, folder := range folders {
		if isEligibleFolder(folder) {
			files, _ := ioutil.ReadDir(computePath(schemaPath, folder.Name()))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".sql") {
					queries := fetchQueries(computePath(schemaPath, folder.Name(), file.Name()))
					scriptPath := computePath(schemaPath, folder.Name(), file.Name())
					scripts = append(scripts, createScript(scriptPath, folder, queries))
				}
			}
		}
	}
	return scripts
}

func createScript(scriptPath string, folder os.FileInfo, queries []Query) Script {
	return Script{Path: scriptPath, Version: folder.Name(), Queries: queries}
}

func isEligibleFolder(folder os.FileInfo) bool {
	return folder.IsDir() && version.Compare(folder.Name(), version.GetCurrentVersion()) == 1
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
		log.Panicf("Error when openning file : %s \n", err.Error())
	}
	var queries []Query
	for _, query := range strings.Split(string(content), ";") {
		queries = append(queries, Query(query))
	}
	return queries
}

func executeScripts(scripts []Script) {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", "db_versioning_test")
	err := db.Connect()
	if err != nil {
		log.Panicf("Connection failed : %s \n", err.Error())
	}
	for _, script := range scripts {
		for _, query := range script.Queries {
			if !query.isEmpty() {
				_, _, err := db.Query(query.GetContent())
				if err != nil {
					upgradeDBVersion(script.Version, script.Path, fmt.Sprintf("failed : %s", err.Error()), db)
					log.Panicf("Error while executing script : %s \n", err.Error())
				}
			}
		}
		upgradeDBVersion(script.Version, script.Path, "ok", db)
	}
	db.Close()
}

func upgradeDBVersion(toVersion, scriptName string, state string, db mysql.Conn) {
	statement, err := db.Prepare("insert into db_version (script, version, state) values (?, ?, ?)")
	if err != nil {
		log.Panicf("Error while preparing the update db_version : %s \n", err.Error())
	}

	_, err = statement.Run(scriptName, toVersion, state)
	if err != nil {
		log.Panicf("Error while updating db_version : %s \n", err.Error())
	}
}

package version

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Version struct {
	VersionNumbers []VersionNumber
}

func (version Version) Compare(versionToCompare Version) int {
	for index, versionNumber := range version.VersionNumbers {
		if versionNumber.isGreaterThan(versionToCompare.VersionNumbers[index]) {
			return 1
		} else if versionNumber.isLowerThan(versionToCompare.VersionNumbers[index]) {
			return -1
		}
	}
	return 0
}

type VersionNumber string

func (versionNumber VersionNumber) isGreaterThan(version VersionNumber) bool {
	return versionNumber > version
}

func (versionNumber VersionNumber) isLowerThan(version VersionNumber) bool {
	return versionNumber < version
}

func Compare(firstVersion string, secondVersion string) int {
	validateVersions(firstVersion, secondVersion)
	firstSplittedVersion, secondSplittedVersion := split(firstVersion, secondVersion)

	return firstSplittedVersion.Compare(secondSplittedVersion)
}

func DisplayCurrentVersion(schema string) {
	fmt.Printf("Current version : %s \n", GetCurrentVersion(schema))
}

func GetCurrentVersion(schema string) string {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "test", "test", schema)
	err := db.Connect()
	if err != nil {
		log.Panicf("Connection failed : %s \n", err.Error())
	}
	versionRow, _, err := db.QueryFirst("select version from db_version order by id desc limit 1")
	if err != nil {
		log.Panicf("Query failed : %s \n", err.Error())
	}
	return versionRow.Str(0)
}

func split(firstVersion string, secondVersion string) (Version, Version) {
	return convertToVersionNumbers(strings.Split(firstVersion, ".")), convertToVersionNumbers(strings.Split(secondVersion, "."))
}

func convertToVersionNumbers(version []string) Version {
	var versionNumber []VersionNumber
	for _, number := range version {
		versionNumber = append(versionNumber, VersionNumber(number))
	}
	return Version{versionNumber}
}

func validateVersions(versions ...string) {
	for _, version := range versions {
		if isFormatValid(version) {
			log.Panicf("Error incompatible version format : %s \n", version)
		}
	}
}

func isFormatValid(version string) bool {
	valid, err := regexp.MatchString("^\\d.\\d.\\d$", version)
	return !valid || err != nil
}

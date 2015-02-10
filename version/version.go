package version

import (
	"db_versioning/db"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Version struct {
	VersionNumbers []VersionNumber
}

type VersionNumber string

func (version Version) Compare(versionToCompare Version) int {
	if version.isGreaterThan(versionToCompare) {
		return 1
	} else if version.IsLowerThan(versionToCompare) {
		return -1
	}
	return 0
}

func Compare(firstVersion string, secondVersion string) int {
	validateVersions(firstVersion, secondVersion)
	firstSplittedVersion, secondSplittedVersion := ConvertToVersionNumbers(firstVersion), ConvertToVersionNumbers(secondVersion)
	return firstSplittedVersion.Compare(secondSplittedVersion)
}

func DisplayCurrentVersion(schema string) {
	fmt.Printf("Current version : %s \n", GetCurrentVersion(schema))
}

func GetCurrentVersion(schema string) string {
	connection := mysql.New("tcp", "", fmt.Sprintf("%s:%d", db.Host, 3306), "test", "test", schema)
	err := connection.Connect()
	if err != nil {
		fmt.Printf("Connection failed : %s \n", err.Error())
		os.Exit(1)
	}
	versionRow, _, err := connection.QueryFirst("select version from db_version where state = 'ok' order by id desc limit 1")
	if err != nil {
		log.Panicf("Query failed : %s \n", err.Error())
	}
	return versionRow.Str(0)
}

func Sort(versions []string) []string {
	sort.Sort(VersionSort(versions))
	return versions
}

type VersionSort []string

func (versionSort VersionSort) Less(i, j int) bool {
	firstVersion, secondVersion := ConvertToVersionNumbers(versionSort[i]), ConvertToVersionNumbers(versionSort[j])
	return firstVersion.IsLowerThan(secondVersion)
}

func (versionSort VersionSort) Swap(i, j int) {
	versionSort[i], versionSort[j] = versionSort[j], versionSort[i]
}

func (versionSort VersionSort) Len() int {
	return len(versionSort)
}

func (version Version) isGreaterThan(versionToCompare Version) bool {
	for index, currentVersion := range version.VersionNumbers {
		if currentVersion.isGreaterThan(versionToCompare.VersionNumbers[index]) {
			return true
		} else if currentVersion.isLowerThan(versionToCompare.VersionNumbers[index]) {
			return false
		}
	}
	return false
}

func (version Version) IsLowerThan(versionToCompare Version) bool {
	for index, currentVersion := range version.VersionNumbers {
		if currentVersion.isLowerThan(versionToCompare.VersionNumbers[index]) {
			return true
		} else if currentVersion.isGreaterThan(versionToCompare.VersionNumbers[index]) {
			return false
		}
	}
	return false
}

func (versionNumber VersionNumber) isGreaterThan(version VersionNumber) bool {
	return convert(versionNumber) > convert(version)
}

func (versionNumber VersionNumber) isLowerThan(version VersionNumber) bool {
	return convert(versionNumber) < convert(version)
}

func convert(version VersionNumber) int {
	convertedVersion, err := strconv.Atoi(string(version))
	if err != nil {
		log.Panicf("Error in conversion of %s : %s", version, err.Error())
	}
	return convertedVersion
}

func ConvertToVersionNumbers(version string) Version {
	var versionNumber []VersionNumber
	for _, number := range strings.Split(version, ".") {
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
	valid, err := regexp.MatchString("^\\d+.\\d+.\\d+$", version)
	return !valid || err != nil
}

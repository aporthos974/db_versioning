package migration

import (
	"db_versioning/version"
	"os"
)

type FolderSort []os.FileInfo

type ScriptSort []Script

func (folderSort FolderSort) Less(i, j int) bool {
	firstFolderSort, secondFolderSort := version.ConvertToVersionNumbers(folderSort[i].Name()), version.ConvertToVersionNumbers(folderSort[j].Name())
	return firstFolderSort.IsLowerThan(secondFolderSort)
}

func (folderSort FolderSort) Swap(i, j int) {
	folderSort[i], folderSort[j] = folderSort[j], folderSort[i]
}

func (folderSort FolderSort) Len() int {
	return len(folderSort)
}

func (scriptSort ScriptSort) Less(i, j int) bool {
	firstFolderSort, secondFolderSort := version.ConvertToVersionNumbers(scriptSort[i].Version), version.ConvertToVersionNumbers(scriptSort[j].Version)
	return firstFolderSort.IsLowerThan(secondFolderSort)
}

func (scriptSort ScriptSort) Swap(i, j int) {
	scriptSort[i], scriptSort[j] = scriptSort[j], scriptSort[i]
}

func (scriptSort ScriptSort) Len() int {
	return len(scriptSort)
}

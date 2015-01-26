package main

import (
	"db_versioning/version"
	"fmt"
)

func main() {
	fmt.Printf("current version : %s\n", version.GetCurrentVersion())
}

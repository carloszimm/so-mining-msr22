package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dlclark/regexp2"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func WriteFolder(basePath string) func(string) {
	return func(folderPath string) {
		err := os.MkdirAll(filepath.Join(basePath, folderPath), os.ModePerm)
		CheckError(err)
	}
}

func RemoveAllFolders(basePath string) func(string) {
	return func(folderPath string) {
		err := os.RemoveAll(filepath.Join(basePath, folderPath))
		CheckError(err)
	}
}

func Regexp2FindAllString(re *regexp2.Regexp, s string) [][]regexp2.Group {
	var matches [][]regexp2.Group
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.Groups())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

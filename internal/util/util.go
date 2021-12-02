package util

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/dlclark/regexp2"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func WriteFolder(folderPath string) {
	err := os.MkdirAll(folderPath, os.ModePerm)
	CheckError(err)
}

func RemoveAllFolders(folderPath string) {
	err := os.RemoveAll(folderPath)
	CheckError(err)
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

func WriteJSON(path string, data interface{}) {
	j, err := json.Marshal(data)
	CheckError(err)

	err = os.WriteFile(path+".json", j, 0644)
	CheckError(err)
}

func CountWords(corpus []string) int {
	count := 0
	for _, w := range corpus {
		count += len(strings.Fields(w))
	}
	return count
}

package util

import (
	"log"

	"github.com/dlclark/regexp2"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
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

package processing

import (
	"regexp"
	"strings"
)

type void struct{}

var member void

func wordPresence(s string) map[string]void {
	wordCount := make(map[string]void)
	for _, word := range strings.Fields(s) {
		wordCount[word] = member
	}

	return wordCount
}

func RmExtremes(corpus []string) []string {
	wordCount := make(map[string][]int)
	var presence map[string]void
	for i, words := range corpus {
		presence = wordPresence(words)
		for key := range presence {
			wordCount[key] = append(wordCount[key], i)
		}
	}
	unCommon := 20
	common := len(corpus) / 2 //50% of documents
	spacesReg := regexp.MustCompile(`\s+`)

	for word, presences := range wordCount {
		length := len(presences)
		if length < unCommon || length > common {
			for _, index := range presences {
				corpus[index] =
					strings.TrimSpace(
						spacesReg.ReplaceAllString(
							strings.ReplaceAll(corpus[index], word, " "), " "))
			}
		}
	}

	return corpus
}

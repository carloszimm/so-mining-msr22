package processing

import (
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

func filter(corpus []string, index int, word string) {
	var str strings.Builder
	for _, w := range strings.Fields(corpus[index]) {
		if w != word {
			str.WriteString(w + " ")
		}
	}
	corpus[index] = strings.TrimSpace(str.String())
}

func RmExtremes(corpus []string) {
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

	for word, presences := range wordCount {
		length := len(presences)
		if length < unCommon || length > common {
			for _, index := range presences {
				filter(corpus, index, word)
			}
		}
	}
}

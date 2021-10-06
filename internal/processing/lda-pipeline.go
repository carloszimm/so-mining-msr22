package processing

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/kljensen/snowball"
)

var (
	specReg = regexp.MustCompile(`[^\w\s']+|(\w)*\d(\w)*`)
	urlReg  = regexp.
		MustCompile(`(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+[\w\-\.,@?^=%&amp;:\/~‌​\+#]*[\w\-\@?^=%&amp‌​;\/~\+#]`)
)

func init() {
	stopwords.LoadStopWordsFromFile(filepath.Join("assets", "stopwords.txt"), "en", "\n")
}

func SetupLDAPipeline(posts []*types.Post, field string) <-chan string {
	out := processPosts(posts, field)
	out = removeTags(out)
	out = removeURLs(out)
	out = removeSpecialChars(out)
	out = removeStopWords(out)
	out = stem(out)
	return out
}

func processPosts(posts []*types.Post, field string) <-chan string {
	out := make(chan string)
	go func() {
		for _, val := range posts {
			out <- types.GetFieldString(val, field)
		}
		close(out)
	}()
	return out
}

func removeTags(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			// load the HTML document
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
			util.CheckError(err)

			// remove code, pre, and blockquotes tags
			doc.Find("code, pre, blockquotes").Remove()

			// get text from html
			out <- doc.Find("body").Text()
		}
		close(out)
	}()
	return out
}

func removeURLs(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			out <- urlReg.ReplaceAllString(text, " ")
		}
		close(out)
	}()
	return out
}

func removeSpecialChars(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			// remove strange chars, punctuations, and numbers
			out <- specReg.ReplaceAllString(text, " ")
		}
		close(out)
	}()
	return out
}

// already maps to lowercase letters
// https://github.com/bbalet/stopwords/blob/master/stopwords.go
func removeStopWords(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			out <- stopwords.CleanString(text, "en", false)
		}
		close(out)
	}()
	return out
}

// already maps to lowercase letters
// https://github.com/kljensen/snowball/blob/master/english/stem.go
func stem(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			textSplitted := strings.Fields(text)
			text = ""

			for _, textPart := range textSplitted {
				stemmed, err := snowball.Stem(textPart, "english", false)
				util.CheckError(err)
				text += (stemmed + " ")
			}

			out <- strings.TrimSpace(text)
		}
		close(out)
	}()
	return out
}

package processing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/carloszimm/stack-mining/types"
	"github.com/carloszimm/stack-mining/util"
	"github.com/kljensen/snowball"
)

var (
	puctReg  *regexp.Regexp = regexp.MustCompile(`\p{P}`)
	spaceReg                = regexp.MustCompile(`\s+`)
	specReg                 = regexp.MustCompile(`[^\w\s\p{P}]`)
)

func init() {
	stopwords.LoadStopWordsFromFile("stopwords.txt", "en", "\n")
}

func SetupPipeline(posts []*types.Post, field string) <-chan string {
	out := processPosts(posts, field)
	out = removeTags(out)
	out = removeStrangeChars(out)
	out = removeStopWords(out)
	out = removePunct(out)
	out = stem(out)
	return out
}

func processPosts(posts []*types.Post, field string) <-chan string {
	out := make(chan string)
	go func() {
		for _, val := range posts {
			out <- util.GetFieldString(val, field)
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

func removeStrangeChars(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			// remove strange chars
			out <- specReg.ReplaceAllString(text, "")
		}
		close(out)
	}()
	return out
}

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

func removePunct(in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		for text := range in {
			// remove puctuation and extra space
			//out <- strings.TrimSpace(spaceReg.ReplaceAllString(puctReg.ReplaceAllString(text, ""), " "))
			out <- puctReg.ReplaceAllString(text, "")
		}
		close(out)
	}()
	return out
}

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
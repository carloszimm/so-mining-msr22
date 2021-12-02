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
	"mvdan.cc/xurls/v2"
)

var (
	regularCharsReg      = regexp.MustCompile(`[^\p{Common}\p{Latin}]+`)
	spacesReg            = regexp.MustCompile(`\s+`)
	pointSymbolNumberReg = regexp.MustCompile(`[\._$\-@\d]+`)
	puctReg              = regexp.MustCompile(`[\p{P}\p{S}\p{Z}]+`) //puctuation, symbol, separator
)

func init() {
	stopwords.LoadStopWordsFromFile(filepath.Join("assets", "stopwords.txt"), "en", "\n")
}

func SetupLDAPipeline(posts []*types.Post, field string) <-chan []string {
	out := processPosts(posts, field)
	out = rmTags(out)
	out = rmStrangeChars(out)
	out = cleanSpaces(out)
	out = rmURLs(out)
	out = rmPointSymbolNumber(out)
	out = rmStopWords(out)
	out = rmPunct(out)
	out = stem(out)
	return filterExtremes(out)
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

func rmTags(in <-chan string) <-chan string {
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

func rmStrangeChars(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			// remove strange chars
			out <- strings.ReplaceAll(regularCharsReg.ReplaceAllString(text, " "), "\uFFFD", " ")
		}
		close(out)
	}()
	return out
}

func cleanSpaces(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			out <- strings.TrimSpace(spacesReg.ReplaceAllString(text, " "))
		}
		close(out)
	}()
	return out
}

func rmURLs(in <-chan string) <-chan string {
	out := make(chan string)
	rxStrict := xurls.Strict()
	go func() {
		for text := range in {
			out <- rxStrict.ReplaceAllString(text, " ")
		}
		close(out)
	}()
	return out
}

/*
temporarily removed given problems with some mappings like:
 -so -> shared-library
 - if -> if-statement
func applyJargon(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			var newText strings.Builder

			stream := jargon.TokenizeString(text).Filter(stackoverflow.Tags)
			for stream.Scan() {
				token := stream.Token()
				if token.IsLemma() {
					newText.WriteString(
						strings.ReplaceAll(
							strings.ReplaceAll(token.String(), "#", "sharp"), "+", "plus"))
				} else {
					newText.WriteString(token.String())
				}
			}
			util.CheckError(stream.Err())

			out <- strings.TrimSpace(newText.String())
		}
		close(out)
	}()
	return out
} */

func rmPointSymbolNumber(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			newText := ""
			for _, txt := range strings.Fields(text) {
				newText += (pointSymbolNumberReg.ReplaceAllString(txt, "") + " ")
			}
			out <- strings.TrimSpace(newText)
		}
		close(out)
	}()
	return out
}

func rmPunct(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for text := range in {
			out <- puctReg.ReplaceAllString(text, " ")
		}
		close(out)
	}()
	return out
}

// already maps to lowercase letters and removes exceding spaces
// https://github.com/bbalet/stopwords/blob/master/stopwords.go
func rmStopWords(in <-chan string) <-chan string {
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
			newText := ""
			for _, textPart := range strings.Fields(text) {
				if len(textPart) > 1 { //skip one-letter words
					stemmed, err := snowball.Stem(textPart, "english", false)
					util.CheckError(err)

					if len(stemmed) > 1 { //skip one-letter words
						newText += (stemmed + " ")
					}
				}
			}

			out <- strings.TrimSpace(newText)
		}
		close(out)
	}()
	return out
}

func filterExtremes(in <-chan string) <-chan []string {
	out := make(chan []string)
	go func() {
		var corpus []string
		for text := range in {
			corpus = append(corpus, text)
		}
		RmExtremes(corpus)
		out <- corpus
		close(out)
	}()
	return out
}

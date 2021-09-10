package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/gocarina/gocsv"
	"github.com/kljensen/snowball"

	"github.com/carloszimm/stack-mining/lda"
	"github.com/carloszimm/stack-mining/types"
	"github.com/carloszimm/stack-mining/util"
)

var filesPath string = filepath.Join("..", "..", "Data Explorer", "9-5-2021", "rxswift")

func init() {
	stopwords.LoadStopWordsFromFile("stopwords.txt", "en", "\n")
}

func main() {
	files, err := ioutil.ReadDir(filesPath)
	util.CheckError(err)

	var count = 0
	c := make(chan []*types.Post)

	for i, f := range files {
		go readCSV(filepath.Join(filesPath, f.Name()), c)
		count = i
	}

	var resultPosts [][]*types.Post

	for i := 0; i < count+1; i++ {
		posts := <-c
		resultPosts = append(resultPosts, posts)
	}

	mergedPosts := mergeArrays(resultPosts)

	//fmt.Println("Total questions:", len(mergedPosts))
	corpus := make([]string, len(mergedPosts))
	for _, val := range mergedPosts {
		corpus = append(corpus, val)
	}

	lda.LDA(corpus)
}

func mergeArrays(postsArray [][]*types.Post) map[int]string {
	result := make(map[int]string)

	for _, posts := range postsArray {
		for _, post := range posts {
			if _, ok := result[post.Id]; !ok {
				result[post.Id] = post.Body
			}
		}
	}

	return result
}

func readCSV(path string, c chan []*types.Post) {
	postsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	util.CheckError(err)
	defer postsFile.Close()

	posts := []*types.Post{}

	if err := gocsv.UnmarshalFile(postsFile, &posts); err != nil {
		panic(err)
	}

	cleanText(&posts)

	c <- posts
}

func cleanText(posts *[]*types.Post) {

	// puctuation
	puctReg := regexp.MustCompile(`\p{P}`)
	spaceReg := regexp.MustCompile(`\s+`)
	specReg := regexp.MustCompile(`[^\w\s\p{P}]`)

	for _, post := range *posts {
		// load the HTML document
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(post.Body))
		util.CheckError(err)

		// remove code, pre, and blockquotes tags
		doc.Find("code, pre, blockquotes").Remove()

		// get text from html
		text := doc.Find("body").Text()

		// remove strange chars
		text = specReg.ReplaceAllString(text, "")

		text = stopwords.CleanString(text, "en", false)

		// remove puctuation and extra space
		text = strings.TrimSpace(spaceReg.ReplaceAllString(puctReg.ReplaceAllString(text, ""), " "))

		text = stem(text)

		post.Body = text
	}
}

func stem(text string) string {
	textSplitted := strings.Split(text, " ")
	text = ""

	for _, textPart := range textSplitted {
		stemmed, err := snowball.Stem(textPart, "english", false)
		util.CheckError(err)
		text += " " + stemmed
	}

	return strings.TrimSpace(text)
}

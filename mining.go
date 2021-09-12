package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/gocarina/gocsv"
	"github.com/kljensen/snowball"

	"github.com/carloszimm/stack-mining/lda"
	"github.com/carloszimm/stack-mining/types"
	"github.com/carloszimm/stack-mining/util"
)

var filesPath string = filepath.Join("data explorer", "rxswift")

func init() {
	stopwords.LoadStopWordsFromFile("stopwords.txt", "en", "\n")
}

func main() {
	files, err := ioutil.ReadDir(filesPath)
	util.CheckError(err)

	posts := readCSVs(files)

	//fmt.Println("Total questions:", len(mergedPosts))
	corpus := make([]string, 0, len(posts))
	for _, val := range posts {
		corpus = append(corpus, val.Body)
	}

	lda.LDA(corpus)
}

func readCSVs(files []fs.FileInfo) []*types.Post {
	c := make(chan []*types.Post)

	for _, f := range files {
		go readCSV(filepath.Join(filesPath, f.Name()), c)
	}

	var resultPosts [][]*types.Post

	for i := 0; i < len(files); i++ {
		posts := <-c
		resultPosts = append(resultPosts, posts)
	}

	return sortPosts(removeDuplicates(resultPosts))
}

func sortPosts(posts []*types.Post) []*types.Post {
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].Id < posts[j].Id
	})
	return posts
}

func removeDuplicates(postsArray [][]*types.Post) []*types.Post {
	resultSet := make(map[int]*types.Post)

	for _, posts := range postsArray {
		for _, post := range posts {
			if _, ok := resultSet[post.Id]; !ok {
				resultSet[post.Id] = post
			}
		}
	}

	result := make([]*types.Post, 0, len(resultSet))
	for _, post := range resultSet {
		result = append(result, post)
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

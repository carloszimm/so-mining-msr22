package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/gocarina/gocsv"
	"github.com/reiver/go-porterstemmer"
)

var filesPath string = filepath.Join("..", "..", "Data Explorer", "31-08-2021", "rxjs")

type Post struct {
	Id           int    `csv:"Id"`
	Body         string `csv:"body"`
	CreationDate string `csv:"-"`
}

func init() {
	stopwords.LoadStopWordsFromFile("stopwords.txt", "en", "\n")
}

func main() {
	files, err := ioutil.ReadDir(filesPath)
	if err != nil {
		log.Fatal(err)
	}

	var count = 0
	c := make(chan []*Post)

	for i, f := range files {
		go readCSV(filepath.Join(filesPath, f.Name()), c)
		count = i
	}

	var resultPosts [][]*Post

	for i := 0; i < count+1; i++ {
		posts := <-c
		resultPosts = append(resultPosts, posts)
	}

	mergedPosts := mergeArrays(resultPosts)

	fmt.Println("Total questions:", len(mergedPosts))
}

func mergeArrays(postsArray [][]*Post) map[int]string {
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

func readCSV(path string, c chan []*Post) {
	postsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		panic(err)
	}
	defer postsFile.Close()

	posts := []*Post{}

	if err := gocsv.UnmarshalFile(postsFile, &posts); err != nil {
		panic(err)
	}

	cleanText(&posts)

	c <- posts
}

func cleanText(posts *[]*Post) {

	// puctuation
	puctReg := regexp.MustCompile(`\p{P}`)
	spaceReg := regexp.MustCompile(`\s+`)

	for _, post := range *posts {
		// load the HTML document
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(post.Body))
		if err != nil {
			log.Fatal(err)
		}

		// remove code, pre, and blockquotes tags
		doc.Find("code, pre, blockquotes").Remove()

		// get text from html
		text := doc.Find("body").Text()

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
		text += " " + porterstemmer.StemString(textPart)
	}

	return strings.TrimSpace(text)
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
)

var filesPath string = filepath.Join("..", "..", "Data Explorer", "31-08-2021", "rxjs")

type Post struct {
	Id           int    `csv:"Id"`
	Body         string `csv:"body"`
	CreationDate string `csv:"-"`
}

func main() {
	files, err := ioutil.ReadDir(filesPath)
	if err != nil {
		log.Fatal(err)
	}

	var count = 0
	c := make(chan map[int]string)

	for i, f := range files {
		go readCSV(filepath.Join(filesPath, f.Name()), c)
		count = i
	}

	var resultPosts []map[int]string

	for i := 0; i < count+1; i++ {
		posts := <-c
		resultPosts = append(resultPosts, posts)
	}

	mergedPosts := mergeMaps(resultPosts)

	fmt.Println(len(mergedPosts))
}

func mergeMaps(maps []map[int]string) map[int]string {
	result := make(map[int]string)

	for _, mapPosts := range maps {
		for id, body := range mapPosts {
			if _, ok := result[id]; !ok {
				result[id] = body
			}
		}
	}

	return result
}

func readCSV(path string, c chan map[int]string) {
	postsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		panic(err)
	}
	defer postsFile.Close()

	posts := []*Post{}

	if err := gocsv.UnmarshalFile(postsFile, &posts); err != nil {
		panic(err)
	}

	m := make(map[int]string)

	for _, post := range posts {
		m[post.Id] = post.Body
	}

	c <- m
}

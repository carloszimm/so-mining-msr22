package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/carloszimm/stack-mining/csv"
	"github.com/carloszimm/stack-mining/lda"
	"github.com/carloszimm/stack-mining/processing"
	"github.com/carloszimm/stack-mining/util"
)

var filesPath string = filepath.Join("data explorer", "rxswift")

func main() {
	files, err := ioutil.ReadDir(filesPath)
	util.CheckError(err)

	posts := csv.ReadPostsCSVs(filesPath, files)

	//fmt.Println("Total questions:", len(mergedPosts))
	corpus := make([]string, 0, len(posts))

	out := processing.SetupPipeline(posts, "Body")
	for text := range out {
		corpus = append(corpus, text)
	}

	lda.LDA(corpus)
}

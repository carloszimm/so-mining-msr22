package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/carloszimm/stack-mining/config"
	"github.com/carloszimm/stack-mining/csv"
	"github.com/carloszimm/stack-mining/lda"
	"github.com/carloszimm/stack-mining/processing"
	"github.com/carloszimm/stack-mining/util"
)

var configs []config.Config

func init() {
	configs = config.ReadConfig()
}

func main() {
	var filesPath string

	for _, cfg := range configs {
		filesPath = filepath.Join(config.DataExplorerPath, cfg.Dir)

		files, err := ioutil.ReadDir(filesPath)
		util.CheckError(err)

		posts := csv.ReadPostsCSVs(filesPath, files)
		log.Println("Processing:", len(posts), "documents for", cfg.Dir, "using", cfg.Field, "field")

		corpus := make([]string, 0, len(posts))

		out := processing.SetupPipeline(posts, cfg.Field)
		for text := range out {
			corpus = append(corpus, text)
		}
		log.Println("Preprocessing finished! Running LDA...")

		if cfg.MinTopics > 0 {
			for i := cfg.MinTopics; i <= cfg.MaxTopics; i++ {
				docTopicDist, topicWordDist := lda.LDA(i, cfg.SampleWords, corpus)
				fmt.Println(docTopicDist, topicWordDist)
			}
		}

		log.Println("LDA finished!")
	}
}

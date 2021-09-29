package main

import (
	"log"
	"path/filepath"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtils "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/lda"
	"github.com/carloszimm/stack-mining/internal/processing"
	"github.com/carloszimm/stack-mining/internal/types"
)

var configs []config.Config

func init() {
	configs = config.ReadConfig()
}

func main() {
	var filesPath string

	for _, cfg := range configs {
		filesPath = filepath.Join(config.CONSOLIDATED_SOURCES_PATH, cfg.FileName+".csv")

		c := make(chan []*types.Post)
		go csvUtils.ReadPostsCSV(filesPath, c)
		posts := <-c

		log.Println("Processing:", len(posts), "documents for", cfg.FileName, "using", cfg.Field, "field")

		corpus := make([]string, 0, len(posts))

		out := processing.SetupPipeline(posts, cfg.Field)
		for text := range out {
			corpus = append(corpus, text)
		}
		log.Println("Preprocessing finished!")
		log.Println("Running LDA...")

		if cfg.MinTopics > 0 {
			csvUtils.WriteFolder(cfg.FileName)
			for i := cfg.MinTopics; i <= cfg.MaxTopics; i++ {
				log.Println("Running for", i, "topics")
				docTopicDist, topicWordDist := lda.LDA(i, cfg.SampleWords, corpus)

				go csvUtils.WriteTopicDist(cfg, i, topicWordDist)
				go csvUtils.WriteDocTopicDist(cfg, i, posts, docTopicDist)
			}
		}

		log.Println("LDA finished!")
	}
}

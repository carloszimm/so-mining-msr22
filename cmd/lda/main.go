package main

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtils "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/lda"
	"github.com/carloszimm/stack-mining/internal/processing"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/golang-module/carbon/v2"
)

var configs []config.Config

func init() {
	configs = config.ReadConfig()
}

func combineTitleBody(posts []*types.Post) {
	for _, post := range posts {
		if post.PostTypeId == 2 {
			parentPost := types.SearchPost(posts, post.ParentId)
			post.Body = parentPost.Title + " " + post.Body
		} else {
			post.Body = post.Title + " " + post.Body
		}
	}
}

func main() {
	var filesPath string

	for _, cfg := range configs {
		filesPath = filepath.Join(config.CONSOLIDATED_SOURCES_PATH, cfg.FileName+".csv")

		c := make(chan []*types.Post)
		go csvUtils.ReadPostsCSV(filesPath, c)
		posts := <-c

		if cfg.CombineTitleBody {
			combineTitleBody(posts)
		}

		log.Println("Processing:", len(posts), "documents for",
			cfg.FileName+".csv", "using", cfg.Field, "field")

		out := processing.SetupLDAPipeline(posts, cfg.Field)
		corpus := <-out

		log.Println("Preprocessing finished!")
		totalWords := util.CountWords(corpus)
		log.Println("Total words:", totalWords)

		log.Println("Running LDA...")

		if cfg.MinTopics > 0 {
			var perplexities []float64

			basePath := filepath.Join(config.LDA_RESULT_PATH,
				strings.ReplaceAll(carbon.Now().ToDateTimeString(), ":", "-"))

			sum := types.Summary{
				MinTopics:  cfg.MinTopics,
				MaxTopics:  cfg.MaxTopics,
				TotalWords: totalWords,
				TotalDocs:  len(posts),
				StartTime:  carbon.Now().ToDayDateTimeString()}

			for i := cfg.MinTopics; i <= cfg.MaxTopics; i++ {
				log.Println("Running for", i, "topics")
				docTopicDist, topicWordDist, perplexity := lda.LDA(i, corpus)

				perplexities = append(perplexities, perplexity)

				basePathWithTopic := filepath.Join(basePath, strconv.Itoa(i))
				//create folder
				util.WriteFolder(basePathWithTopic)

				//write topic and doc x topic distribution
				var wg sync.WaitGroup
				wg.Add(2)
				go csvUtils.WriteTopicDist(&wg, cfg, basePathWithTopic, i, topicWordDist)
				go csvUtils.WriteDocTopicDist(&wg, cfg, basePathWithTopic, i, posts, docTopicDist)
				wg.Wait()
			}

			csvUtils.WritePerplexities(cfg, basePath, perplexities)

			sum.EndTime = carbon.Now().ToDayDateTimeString()
			types.WriteSummary(basePath, sum)
		}

		log.Println("LDA finished!")
	}
}

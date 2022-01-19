package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtil "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

const (
	DATE             = "2022-01-12 02-21-28"
	NUMBER_OF_TOPICS = "23"
)

var (
	DOCTOPICS_PATH       = path.Join(config.LDA_RESULT_PATH, DATE, NUMBER_OF_TOPICS, fmt.Sprintf("all_withAnswers_doctopicdist_%s_Body.csv", NUMBER_OF_TOPICS))
	POSTS_PATH           = path.Join(config.CONSOLIDATED_SOURCES_PATH, "all_withAnswers.csv")
	EXTRACTED_POSTS_PATH = path.Join("assets", "extracted-posts")
)

// topic to be searched for
const TOPIC = 22

func main() {
	util.WriteFolder(EXTRACTED_POSTS_PATH)

	// loads docXtopic csv
	docTopics := csvUtil.ReadDocTopic(DOCTOPICS_PATH)
	shares := types.GetTopicShare(docTopics)

	var posts []int
	for _, s := range shares[TOPIC] {
		posts = append(posts, s.PostId)
	}
	postsBytes, _ := json.Marshal(posts)

	err := os.WriteFile(path.Join(EXTRACTED_POSTS_PATH, fmt.Sprintf("%d.json", TOPIC)), postsBytes, 0644)
	util.CheckError(err)
}

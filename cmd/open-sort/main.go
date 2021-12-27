package main

import (
	"math/rand"
	"path"
	"time"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtil "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

const NUMBER_POSTS = 15

var (
	DOCTOPICS_PATH = path.Join(config.LDA_RESULT_PATH, "2021-12-21 00-01-48", "25", "all_withAnswers_doctopicdist_25_Body.csv")
	POSTS_PATH     = path.Join(config.CONSOLIDATED_SOURCES_PATH, "all_withAnswers.csv")
)

var random *rand.Rand

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	random = rand.New(s1)
}

func getRandomPostIndex(length int) int {
	return random.Intn(length)
}

func main() {

	docTopics := csvUtil.ReadDocTopic(DOCTOPICS_PATH)

	shares := types.GetTopicShare(docTopics)
	randomPosts := make(map[int][]int)

	for i := 0; i < len(shares); i++ {
		seemIndexes := make(map[int]bool)
		for len(randomPosts[i]) < NUMBER_POSTS {
			var randomPostIndex int
			for {
				randomPostIndex = getRandomPostIndex(len(shares[i]))
				if shares[i][randomPostIndex].Proportion > 0.1 {
					break
				}
			}
			if _, ok := seemIndexes[randomPostIndex]; !ok {
				randomPosts[i] = append(randomPosts[i], shares[i][randomPostIndex].PostId)
				seemIndexes[randomPostIndex] = true
			}
		}
	}

	util.WriteFolder(config.OPENSORT_RESULT_PATH)

	csvUtil.WriteOpenSort(
		path.Join(config.OPENSORT_RESULT_PATH, "opensort.csv"), randomPosts)
}

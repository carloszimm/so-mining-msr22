package csvUtil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

func WriteTopicDist(wg *sync.WaitGroup, cfg config.Config, topics int, data [][]types.WordDist) {
	defer wg.Done()

	filePath := filepath.Join(config.LDA_RESULT_PATH,
		cfg.FileName, cfg.Field, strconv.Itoa(topics),
		fmt.Sprintf("%s_%s_%d_%s.csv", cfg.FileName, "topicdist", topics, cfg.Field))

	file, err := os.Create(filePath)
	util.CheckError(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Topics", "Words"})
	util.CheckError(err)

	for topic, record := range data {
		words := ""
		for _, wordDist := range record {
			words += (wordDist.Word + " ")
		}
		words = strings.TrimSpace(words)
		err := writer.Write([]string{strconv.Itoa(topic), words})
		util.CheckError(err)
	}
}

func WriteDocTopicDist(wg *sync.WaitGroup, cfg config.Config, topics int, posts []*types.Post, data [][]types.TopicDist) {
	defer wg.Done()

	filePath := filepath.Join(config.LDA_RESULT_PATH,
		cfg.FileName, cfg.Field, strconv.Itoa(topics),
		fmt.Sprintf("%s_%s_%d_%s.csv", cfg.FileName, "doctopicdist", topics, cfg.Field))

	file, err := os.Create(filePath)
	util.CheckError(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Post_ID", "Dominant_Topic", "Topic_Percent"})
	util.CheckError(err)

	for doc, topics := range data {
		err := writer.Write([]string{strconv.Itoa(posts[doc].Id),
			strconv.Itoa(topics[0].Topic), fmt.Sprintf("%.3f", topics[0].Probability*100)})
		util.CheckError(err)
	}
}

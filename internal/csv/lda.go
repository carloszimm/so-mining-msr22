package csvUtil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

var regex = regexp.MustCompile(`\[|\]`)

func WriteTopicDist(wg *sync.WaitGroup, cfg config.Config, topics int, data [][]types.WordDist) {
	defer wg.Done()

	baseFilePath := filepath.Join(config.LDA_RESULT_PATH,
		cfg.FileName, cfg.Field, strconv.Itoa(topics),
		fmt.Sprintf("%s_%s_%d_%s", cfg.FileName, "topicdist", topics, cfg.Field))

	var wGroup sync.WaitGroup
	wGroup.Add(2)
	func() {
		defer wGroup.Done()
		writeTopWords(baseFilePath, cfg, topics, data)
	}()
	func() {
		defer wGroup.Done()
		writeComplete(baseFilePath, cfg, topics, data)
	}()
	wGroup.Wait()
}
func writeTopWords(baseFilePath string, cfg config.Config, topics int, data [][]types.WordDist) {
	filePath := baseFilePath + " - topwords.csv"

	file, err := os.Create(filePath)
	util.CheckError(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Topics", fmt.Sprintf("%s_%d_%s", "Top", cfg.SampleWords, "Words")})
	util.CheckError(err)

	for topic, record := range data {
		topWordsSlice := record[:cfg.SampleWords]
		topWords := ""
		for _, word := range topWordsSlice {
			topWords += (word.Word + " ")
		}

		err := writer.Write([]string{strconv.Itoa(topic), strings.TrimSpace(topWords)})
		util.CheckError(err)
	}
}
func writeComplete(baseFilePath string, cfg config.Config, topics int, data [][]types.WordDist) {
	filePath := baseFilePath + ".csv"

	file, err := os.Create(filePath)
	util.CheckError(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Topics", "Words"})
	util.CheckError(err)

	for topic, record := range data {
		words := regex.ReplaceAllString(fmt.Sprint(record), "")

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

	err = writer.Write([]string{"Post_ID", "Dominant_Topic", "Topic_Proportion", "Topics"})
	util.CheckError(err)

	for doc, topics := range data {
		allTopics := regex.ReplaceAllString(fmt.Sprint(topics), "")
		err := writer.Write([]string{strconv.Itoa(posts[doc].Id),
			strconv.Itoa(topics[0].Topic), fmt.Sprintf("%f", topics[0].Probability),
			allTopics})
		util.CheckError(err)
	}
}

func WritePerplexities(cfg config.Config, perplexities []float64) {
	filePath := filepath.Join(config.LDA_RESULT_PATH,
		cfg.FileName, cfg.Field, "perplexity.csv")

	file, err := os.Create(filePath)
	util.CheckError(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Topics", "Perplexity", "Decrease"})
	util.CheckError(err)

	originalNumber := perplexities[0]

	for i, perplexity := range perplexities {
		decrease := ((originalNumber - perplexity) / originalNumber) * 100

		err := writer.Write([]string{strconv.Itoa(cfg.MinTopics + i),
			fmt.Sprint(perplexity), fmt.Sprint(decrease)})
		util.CheckError(err)

		originalNumber = perplexity
	}
}

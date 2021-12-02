package types

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/carloszimm/stack-mining/internal/util"
)

type Summary struct {
	StartTime  string
	EndTime    string
	MinTopics  int
	MaxTopics  int
	TotalDocs  int
	TotalWords int
}

func WriteSummary(path string, sum Summary) {
	template := "Start Time: %v\nEnd Time: %v\nTopics Range: %v-%v\nTotal of Documents: %v\nTotal of Words: %v\n"
	text := fmt.Sprintf(template, sum.StartTime, sum.EndTime, sum.MinTopics, sum.MaxTopics, sum.TotalDocs, sum.TotalWords)
	err := os.WriteFile(filepath.Join(path, "summary.txt"), []byte(text), 0644)
	util.CheckError(err)
}

type TopicDist struct {
	Topic       int
	Probability float64
}

func (td TopicDist) String() string {
	return fmt.Sprintf("%v(%f)", td.Topic, td.Probability)
}

type WordDist struct {
	Word        string
	Probability float64
}

func (wd WordDist) String() string {
	return fmt.Sprintf("%v(%f)", wd.Word, wd.Probability)
}

func SortLdaDesc(i interface{}) {
	switch dist := i.(type) {
	case []TopicDist:
		sort.SliceStable(dist, func(i, j int) bool {
			return dist[i].Probability > dist[j].Probability
		})
	case []WordDist:
		sort.SliceStable(dist, func(i, j int) bool {
			return dist[i].Probability > dist[j].Probability
		})
	}
}

package csvUtil

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/gocarina/gocsv"
)

func ReadDocTopic(path string) []*types.DocTopic {
	docTopicFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	util.CheckError(err)
	defer docTopicFile.Close()

	docTopics := []*types.DocTopic{}

	err = gocsv.UnmarshalFile(docTopicFile, &docTopics)
	util.CheckError(err)

	return docTopics
}

func WriteOpenSort(path string, randomPosts map[int][]int) {
	f, err := os.Create(path)
	util.CheckError(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for i := 0; i < len(randomPosts); i++ {
		var bodies []string
		bodies = append(bodies, fmt.Sprint(i))
		for _, randomPostId := range randomPosts[i] {
			bodies = append(bodies,
				fmt.Sprintf(`=HYPERLINK("%s%d", "%d")`,
					"https://stackoverflow.com/questions/", randomPostId, randomPostId))
		}
		err = w.Write(bodies)
		util.CheckError(err)
	}
}

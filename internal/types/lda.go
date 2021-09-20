package types

import (
	"fmt"
	"sort"
)

type TopicDist struct {
	Topic       int
	Probability float64
}

type WordDist struct {
	Word        string
	Probability float64
}

func (wd WordDist) String() string {
	return fmt.Sprintf("%v(%.3f%%)", wd.Word, wd.Probability*100)
}

func SortLdaDesc(i interface{}) {
	switch dist := i.(type) {
	case []TopicDist:
	case []WordDist:
		sort.SliceStable(dist, func(i, j int) bool {
			return dist[i].Probability > dist[j].Probability
		})
	}
}

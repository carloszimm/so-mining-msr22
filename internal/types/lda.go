package types

import (
	"fmt"
	"sort"
)

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

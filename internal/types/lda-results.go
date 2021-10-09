package types

type DocTopic struct {
	PostId        int     `csv:"Post_ID"`
	DominantTopic int     `csv:"Dominant_Topic"`
	Proportion    float64 `csv:"Topic_Proportion"`
	Topics        string  `csv:"-"`
}

func GetTopicShare(doctopics []*DocTopic) map[int][]*DocTopic {
	topicShare := make(map[int][]*DocTopic)
	for _, doctopic := range doctopics {
		topicShare[doctopic.DominantTopic] =
			append(topicShare[doctopic.DominantTopic], doctopic)
	}

	return topicShare
}

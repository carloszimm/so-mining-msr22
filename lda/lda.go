package lda

import (
	"github.com/carloszimm/stack-mining/types"
	"github.com/carloszimm/stack-mining/util"
	"github.com/james-bowman/nlp"
)

func LDA(topics int, corpus []string) (*[][]types.TopicDist, *[][]types.WordDist) {

	vectoriser := nlp.NewCountVectoriser()
	lda := nlp.NewLatentDirichletAllocation(topics)
	pipeline := nlp.NewPipeline(vectoriser, lda)

	docsOverTopics, err := pipeline.FitTransform(corpus...)
	util.CheckError(err)

	// Examine Document over topic probability distribution
	dr, dc := docsOverTopics.Dims()

	docTopicDist := make([][]types.TopicDist, dc)

	for doc := 0; doc < dc; doc++ {
		docTopicDist[doc] = make([]types.TopicDist, 0, dr)

		for topic := 0; topic < dr; topic++ {
			docTopicDist[doc] = append(docTopicDist[doc],
				types.TopicDist{Topic: topic, Probability: docsOverTopics.At(topic, doc)})
		}
	}

	// Examine Topic over word probability distribution
	topicsOverWords := lda.Components()
	tr, tc := topicsOverWords.Dims()

	vocab := make([]string, len(vectoriser.Vocabulary))
	for k, v := range vectoriser.Vocabulary {
		vocab[v] = k
	}

	topicWordDist := make([][]types.WordDist, tr)

	for topic := 0; topic < tr; topic++ {
		topicWordDist[topic] = make([]types.WordDist, 0, tc)

		for word := 0; word < tc; word++ {
			topicWordDist[topic] = append(topicWordDist[topic],
				types.WordDist{Word: vocab[word], Probability: topicsOverWords.At(topic, word)})
		}
	}

	return &docTopicDist, &topicWordDist
}

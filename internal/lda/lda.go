package lda

import (
	"fmt"

	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/james-bowman/nlp"
)

var stopWords = []string{
	"differ", "specif", "deal", "prefer", "easili", "easier", "mind",

	"current", "solv", "proper", "modifi", "explain", "hope", "help", "wonder", "altern", "sens", "entir", "ps",

	"solut", "achiev", "approach", "answer", "requir", "lot", "feel", "pretti", "easi", "goal", "think",
	"complex", "eleg", "improv", "look", "complic", "day",

	"chang", "issu", "add", "edit", "remov", "custom", "suggest", "comment", "ad", "refer", "stackblitz",
	"link", "mention", "detect", "face", "fix", "attach", "perfect", "mark",

	"reason", "suppos", "notic", "snippet", "demo", "line", "piec", "appear",
}

func init() {
	for ch := 'a'; ch <= 'z'; ch++ {
		stopWords = append(stopWords, fmt.Sprintf("%c", ch))
	}
}

func LDA(topics int, corpus []string) ([][]types.TopicDist, [][]types.WordDist, float64) {

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	lda := nlp.NewLatentDirichletAllocation(topics)

	lda.Alpha, lda.Eta = 0.01, 0.01

	vectorisedData, err := vectoriser.FitTransform(corpus...)
	util.CheckError(err)

	docsOverTopics, err := lda.FitTransform(vectorisedData)
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

		types.SortLdaDesc(docTopicDist[doc])
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

		types.SortLdaDesc(topicWordDist[topic])
	}

	return docTopicDist, topicWordDist, lda.Perplexity(vectorisedData)
}

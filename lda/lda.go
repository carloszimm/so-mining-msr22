package lda

import (
	"fmt"
	"os"
	"sort"

	"github.com/carloszimm/stack-mining/util"
	"github.com/james-bowman/nlp"
	"github.com/olekukonko/tablewriter"
)

func LDA(corpus []string) {
	// Create a pipeline with a count vectoriser and LDA transformer for 2 topics
	vectoriser := nlp.NewCountVectoriser()
	lda := nlp.NewLatentDirichletAllocation(20)
	pipeline := nlp.NewPipeline(vectoriser, lda)

	docsOverTopics, err := pipeline.FitTransform(corpus...)
	if err != nil {
		fmt.Printf("Failed to model topics for documents because %v", err)
		return
	}

	// Examine Document over topic probability distribution
	docsOverTopics.Dims()
	//dr, dc := docsOverTopics.Dims()
	/* for doc := 0; doc < dc; doc++ {
		fmt.Printf("\nTopic distribution for document: '%s' -", corpus[doc])
		for topic := 0; topic < dr; topic++ {
			if topic > 0 {
				fmt.Printf(",")
			}
			fmt.Printf(" Topic #%d=%f", topic, docsOverTopics.At(topic, doc))
		}
	} */

	// Examine Topic over word probability distribution
	topicsOverWords := lda.Components()
	tr, tc := topicsOverWords.Dims()

	vocab := make([]string, len(vectoriser.Vocabulary))
	for k, v := range vectoriser.Vocabulary {
		vocab[v] = k
	}
	for topic := 0; topic < tr; topic++ {

		data := make(map[string]float64)

		for word := 0; word < tc; word++ {
			data[vocab[word]] = topicsOverWords.At(topic, word)
		}
		writeData("topic_"+fmt.Sprint(topic), data)
	}
}

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func writeData(topic string, words map[string]float64) {
	p := make(PairList, len(words))

	i := 0
	for k, v := range words {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)

	f, err := os.Create("./lda/results/" + topic + ".txt")
	util.CheckError(err)
	defer f.Close()

	table := tablewriter.NewWriter(f)
	table.SetHeader([]string{"Term", "Probability"})

	for _, k := range p {
		table.Append([]string{k.Key, fmt.Sprintf("%.3f%%", k.Value*100)})
	}

	table.Render()
}

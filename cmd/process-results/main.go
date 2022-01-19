package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtil "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
)

const (
	DATE             = "2022-01-12 02-21-28"
	NUMBER_OF_TOPICS = "23"
)

var (
	DOCTOPICS_PATH         = path.Join(config.LDA_RESULT_PATH, DATE, NUMBER_OF_TOPICS, fmt.Sprintf("all_withAnswers_doctopicdist_%s_Body.csv", NUMBER_OF_TOPICS))
	POSTS_PATH             = path.Join(config.CONSOLIDATED_SOURCES_PATH, "all_withAnswers.csv")
	RESULT_PROCESSING_PATH = path.Join("assets", "result-processing")
)

type topicTotal struct {
	topic      int
	total      int
	percentage float64
}

type popularityRawData struct {
	topic    int
	view     []int
	favorite []int
	score    []int
}

type popularity struct {
	topic       int
	avgView     float64
	avgFavorite float64
	avgScore    float64
}

type topicQuestion struct {
	topic     int
	questions []*types.Post
}

type difficulty struct {
	topic                 int
	percentageNoAccAnswer float64
	medianTime            float64
}

func writeToTxt(path string, header []string, data [][]string) {
	f, err := os.Create(path)
	util.CheckError(err)
	defer f.Close()

	table := tablewriter.NewWriter(f)
	table.SetHeader(header)

	table.AppendBulk(data)

	table.Render()
}

func writeBasicInfo(share []topicTotal) {
	resultPath := path.Join(RESULT_PROCESSING_PATH, "topics_info.txt")
	header := []string{"Topic", "Posts", "Percentage"}
	var results [][]string
	for _, t := range share {
		results = append(results,
			[]string{fmt.Sprint(t.topic), fmt.Sprint(t.total), fmt.Sprintf("%.1f", t.percentage)})
	}
	writeToTxt(resultPath, header, results)
}

func writePopularity(popularities []popularity) {
	resultPath := path.Join(RESULT_PROCESSING_PATH, "popularity.txt")
	header := []string{"Topic", "Avg. View", "Avg. Fav.", "Avg. Score"}
	var results [][]string
	for _, pop := range popularities {
		results = append(results,
			[]string{fmt.Sprint(pop.topic), fmt.Sprintf("%.1f", pop.avgView),
				fmt.Sprintf("%.1f", pop.avgFavorite), fmt.Sprintf("%.1f", pop.avgScore)})
	}
	writeToTxt(resultPath, header, results)
}

func writeDifficulty(difficulties []difficulty) {
	resultPath := path.Join(RESULT_PROCESSING_PATH, "difficulty.txt")
	header := []string{"Topic", "Perc. No Acc. Answer", "Median Time(h)"}
	var results [][]string
	for _, diff := range difficulties {
		results = append(results,
			[]string{fmt.Sprint(diff.topic),
				fmt.Sprintf("%.1f", diff.percentageNoAccAnswer), fmt.Sprintf("%.1f", diff.medianTime)})
	}
	writeToTxt(resultPath, header, results)
}

func basicInfo(shares map[int][]*types.DocTopic) {
	share := make([]topicTotal, 0, len(shares))
	totalDocs := 0
	for t, docTopic := range shares {
		total := len(docTopic)
		totalDocs += total
		share = append(share, topicTotal{topic: t, total: total})
	}

	for i := range share {
		share[i].percentage = (float64(share[i].total) / float64(totalDocs)) * 100
	}

	sort.SliceStable(share, func(i, j int) bool {
		return share[i].total > share[j].total
	})

	//fmt.Println(totalDocs)
	writeBasicInfo(share)
}

func calculatePopularity(shares map[int][]*types.DocTopic, posts []*types.Post) {
	popRaw := make([]popularityRawData, 0, len(shares))

	for t, docTopic := range shares {
		popRawData := popularityRawData{topic: t}
		// goes through the posts of that topic
		for _, dT := range docTopic {
			// retrieve info about post
			post := types.SearchPost(posts, dT.PostId)
			if post != nil && post.PostTypeId == 1 { //question post
				popRawData.view = append(popRawData.view, post.ViewCount)
				popRawData.favorite = append(popRawData.favorite, post.FavoriteCount)
				popRawData.score = append(popRawData.score, post.Score)
			}
		}
		popRaw = append(popRaw, popRawData)
	}
	popularities := make([]popularity, 0, len(shares))
	for _, raw := range popRaw {
		p := popularity{topic: raw.topic}
		p.avgView, _ = stats.Mean(stats.LoadRawData(raw.view))
		p.avgFavorite, _ = stats.Mean(stats.LoadRawData(raw.favorite))
		p.avgScore, _ = stats.Mean(stats.LoadRawData(raw.score))
		popularities = append(popularities, p)
	}
	// sort slice descendingly
	sort.SliceStable(popularities, func(i, j int) bool {
		return popularities[i].avgView > popularities[j].avgView
	})

	writePopularity(popularities)
}

func calculateDifficulty(shares map[int][]*types.DocTopic, posts []*types.Post) {
	topicQuestions := make([]topicQuestion, 0, len(shares))

	for t, docTopic := range shares {
		tQ := topicQuestion{topic: t}
		// goes through the posts of that topic
		for _, dT := range docTopic {
			// retrieve info about post
			post := types.SearchPost(posts, dT.PostId)
			if post != nil && post.PostTypeId == 1 { //question
				tQ.questions = append(tQ.questions, post)
			}
		}
		// stores the questions related to the current topic
		topicQuestions = append(topicQuestions, tQ)
	}

	difficulties := make([]difficulty, 0, len(shares))
	for _, topicQ := range topicQuestions {
		diff := difficulty{topic: topicQ.topic}
		noAnswer := 0
		var postWithAnswers []*types.Post
		for _, q := range topicQ.questions {
			if q.AcceptedAnswerId == 0 {
				noAnswer++
			} else {
				postWithAnswers = append(postWithAnswers, q)
			}
		}

		diff.percentageNoAccAnswer = (float64(noAnswer) / float64(len(topicQ.questions)) * 100)

		durations := make([]float64, 0, len(postWithAnswers))
		for _, pAnswer := range postWithAnswers {
			// retrieve info about accepted answer post
			answer := types.SearchPost(posts, pAnswer.AcceptedAnswerId)
			durations = append(durations, answer.CreationDate.Sub(pAnswer.CreationDate.Time).Hours())
		}
		//data := stats.LoadRawData(durations)
		diff.medianTime, _ = stats.Median(durations)

		difficulties = append(difficulties, diff)
	}
	// sort slice descendingly
	sort.SliceStable(difficulties, func(i, j int) bool {
		return difficulties[i].percentageNoAccAnswer > difficulties[j].percentageNoAccAnswer
	})

	writeDifficulty(difficulties)
}

func main() {
	log.Println("Starting Processing...")

	// loads docXtopic csv
	docTopics := csvUtil.ReadDocTopic(DOCTOPICS_PATH)
	shares := types.GetTopicShare(docTopics)

	// loads posts
	c := make(chan []*types.Post)
	go csvUtil.ReadPostsCSV(POSTS_PATH, c)
	posts := <-c

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		basicInfo(shares)
	}()
	go func() {
		defer wg.Done()
		calculatePopularity(shares, posts)
	}()
	go func() {
		defer wg.Done()
		calculateDifficulty(shares, posts)
	}()

	wg.Wait()
	log.Println("Processing Finished!")
}

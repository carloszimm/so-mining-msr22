package main

import (
	"fmt"
	"path"
	"sort"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtil "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/montanaflynn/stats"
)

const (
	DATE             = "2022-01-12 02-21-28"
	NUMBER_OF_TOPICS = "23"
)

var (
	DOCTOPICS_PATH = path.Join(config.LDA_RESULT_PATH, DATE, NUMBER_OF_TOPICS, fmt.Sprintf("all_withAnswers_doctopicdist_%s_Body.csv", NUMBER_OF_TOPICS))
	POSTS_PATH     = path.Join(config.CONSOLIDATED_SOURCES_PATH, "all_withAnswers.csv")
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

func basicInfo(shares map[int][]*types.DocTopic) {
	share := make([]topicTotal, 0, len(shares))
	totalDocs := 0
	for t, docTopic := range shares {
		total := len(docTopic)
		totalDocs += total
		share = append(share, topicTotal{topic: t, total: total})
	}

	for i := range share {
		share[i].percentage = float64(share[i].total) / float64(totalDocs)
	}

	sort.SliceStable(share, func(i, j int) bool {
		return share[i].total > share[j].total
	})

	fmt.Println(share)
	fmt.Println(totalDocs)
}

func calculatePopularity(result chan []popularity, shares map[int][]*types.DocTopic, posts []*types.Post) {
	popRaw := make([]popularityRawData, 0, len(shares))
	for t, docTopic := range shares {
		popRawData := popularityRawData{topic: t}
		for _, dT := range docTopic {
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

	sort.SliceStable(popularities, func(i, j int) bool {
		return popularities[i].avgView > popularities[j].avgView
	})

	//fmt.Println(popularities)
	result <- popularities
}

func calculateDifficulty(result chan []difficulty, shares map[int][]*types.DocTopic, posts []*types.Post) {
	topicQuestions := make([]topicQuestion, 0, len(shares))

	for t, docTopic := range shares {
		tQ := topicQuestion{topic: t}
		for _, dT := range docTopic {
			post := types.SearchPost(posts, dT.PostId)
			if post != nil && post.PostTypeId == 1 { //question
				tQ.questions = append(tQ.questions, post)
			}
		}
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

		diff.percentageNoAccAnswer = float64(noAnswer) / float64(len(topicQ.questions))

		durations := make([]float64, 0, len(postWithAnswers))
		for _, pAnswer := range postWithAnswers {
			answer := types.SearchPost(posts, pAnswer.AcceptedAnswerId)
			durations = append(durations, answer.CreationDate.Sub(pAnswer.CreationDate.Time).Hours())
		}
		//data := stats.LoadRawData(durations)
		diff.medianTime, _ = stats.Median(durations)

		difficulties = append(difficulties, diff)
	}

	sort.SliceStable(difficulties, func(i, j int) bool {
		return difficulties[i].percentageNoAccAnswer > difficulties[j].percentageNoAccAnswer
	})

	//fmt.Println(difficulties)

	result <- difficulties
}

func main() {
	// loads docxtopic csv
	docTopics := csvUtil.ReadDocTopic(DOCTOPICS_PATH)
	shares := types.GetTopicShare(docTopics)

	// loads posts
	c := make(chan []*types.Post)
	go csvUtil.ReadPostsCSV(POSTS_PATH, c)
	posts := <-c

	popularityChannel := make(chan []popularity)
	difficultyChannel := make(chan []difficulty)

	go calculatePopularity(popularityChannel, shares, posts)
	go calculateDifficulty(difficultyChannel, shares, posts)

	basicInfo(shares)

	popularities := <-popularityChannel
	difficulties := <-difficultyChannel

	fmt.Println(popularities)
	fmt.Println(difficulties)
}

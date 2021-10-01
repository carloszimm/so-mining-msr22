package main

import (
	"fmt"
	"path/filepath"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtils "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/processing"
	"github.com/carloszimm/stack-mining/internal/types"
)

func main() {
	operators := types.CreateOperators("rxjs 7.3.0.json", "rxjs")
	filesPath := filepath.Join(config.CONSOLIDATED_SOURCES_PATH, "rxjs.csv")

	c := make(chan []*types.Post)
	go csvUtils.ReadPostsCSV(filesPath, c)
	posts := <-c

	resultChannel := processing.SetupOpsPipeline(posts, operators)

	result := <-resultChannel

	test := types.OperatorCount{}
	for _, val := range result {
		test.OpName = val[5].OpName
		test.Total += val[5].Total
	}
	fmt.Println(test)
}

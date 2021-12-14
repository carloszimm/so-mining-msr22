package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtils "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/processing"
	"github.com/carloszimm/stack-mining/internal/stats"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

var sources []string

func init() {
	files, err := ioutil.ReadDir(config.CONSOLIDATED_SOURCES_PATH)
	util.CheckError(err)

	// store source files' name
	for _, f := range files {
		if !f.IsDir() && !strings.Contains(f.Name(), "all") {
			sources = append(sources, f.Name())
		}
	}
}

func main() {
	operatorsFiles, err := ioutil.ReadDir(config.OPERATORS_PATH)
	util.CheckError(err)

	for _, opFile := range operatorsFiles {
		for _, source := range sources {
			distFileName := strings.TrimSuffix(source, filepath.Ext(source))
			dist := strings.Split(distFileName, "_")[0]

			// (?i) case-insensitive mode
			// .* any character except line break
			if regexp.MustCompile("(?i)" + dist + ".*").MatchString(opFile.Name()) {
				operators := types.CreateOperators(opFile.Name(), dist)
				filesPath := filepath.Join(config.CONSOLIDATED_SOURCES_PATH, source)

				c := make(chan []*types.Post)
				go csvUtils.ReadPostsCSV(filesPath, c)
				posts := <-c

				resultChannel := processing.SetupOpsPipeline(posts, operators)

				// the result is a map consisting of the post id as the key and Operator Count instances that
				// contain operator count and its name
				result := <-resultChannel

				generateResults(
					strings.TrimSuffix(opFile.Name(), filepath.Ext(opFile.Name()))+"_"+distFileName,
					result)
			}
		}
	}
}

func generateResults(path string, result map[int][]types.OperatorCount) {
	// clean up old files
	removeOldFiles(path)

	opsCount := types.AggregateByOperator(result)
	// statistics according to the operators
	// opsCount is only used here
	resultStats := stats.GenerateOpsStats(opsCount)

	var wg sync.WaitGroup
	wg.Add(2)
	// writes the original result as JSON for future usage
	go func() {
		defer wg.Done()
		util.WriteJSON(filepath.Join(config.OPERATORS_RESULT_PATH, path), result)
	}()
	// writes the statistics to CSV
	go func() {
		defer wg.Done()
		csvUtils.WriteOpsSearchResult(filepath.Join(config.OPERATORS_RESULT_PATH, path), resultStats)
	}()
	wg.Wait()
}

func removeOldFiles(path string) {
	resultFiles, err := ioutil.ReadDir(config.OPERATORS_RESULT_PATH)
	util.CheckError(err)

	for _, resultFile := range resultFiles {
		if strings.Contains(resultFile.Name(), path) {
			util.RemoveAllFolders(filepath.Join(config.OPERATORS_RESULT_PATH, resultFile.Name()))
		}
	}
}

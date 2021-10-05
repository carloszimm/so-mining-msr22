package main

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtils "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

func main() {
	distPosts := make(map[string][]*types.Post)
	distPostsWithAnswers := make(map[string][]*types.Post)

	dirs, err := ioutil.ReadDir(config.DATA_EXPLORER_PATH)
	util.CheckError(err)

	for _, f := range dirs {
		if f.IsDir() && f.Name() != "consolidated sources" {
			files, err := ioutil.ReadDir(filepath.Join(config.DATA_EXPLORER_PATH, f.Name()))
			util.CheckError(err)

			var filesWithNoAnswers []fs.FileInfo
			var filesWithAnswers []fs.FileInfo
			for _, file := range files {
				if !strings.Contains(file.Name(), "withAnswers") {
					filesWithNoAnswers = append(filesWithNoAnswers, file)
				} else {
					filesWithAnswers = append(filesWithAnswers, file)
				}
			}
			if len(filesWithNoAnswers) > 0 {
				distPosts[f.Name()] =
					csvUtils.ReadPostsCSVs(filepath.Join(config.DATA_EXPLORER_PATH, f.Name()), filesWithNoAnswers)
			}
			if len(filesWithAnswers) > 0 {
				distPostsWithAnswers[f.Name()] =
					csvUtils.ReadPostsCSVs(filepath.Join(config.DATA_EXPLORER_PATH, f.Name()), filesWithAnswers)
			}
		}
	}

	var wg sync.WaitGroup
	if len(distPosts) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			writeResults("", distPosts)
		}()
	}
	if len(distPostsWithAnswers) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			writeResults("_withAnswers", distPostsWithAnswers)
		}()
	}
	wg.Wait()
}

func writeResults(suffix string, distPosts map[string][]*types.Post) {
	var allPosts []*types.Post

	var wg sync.WaitGroup
	wg.Add(len(distPosts))

	for dist, posts := range distPosts {
		//avoid problem with closure
		dist := dist
		posts := posts
		go func() {
			defer wg.Done()
			csvUtils.WritePostsCSV(filepath.Join(config.CONSOLIDATED_SOURCES_PATH, dist+suffix+".csv"), posts)
		}()

		allPosts = append(allPosts, posts...)
	}
	wg.Wait()

	csvUtils.WritePostsCSV(
		filepath.Join(config.CONSOLIDATED_SOURCES_PATH, "all"+suffix+".csv"),
		csvUtils.SortPosts(allPosts))
}

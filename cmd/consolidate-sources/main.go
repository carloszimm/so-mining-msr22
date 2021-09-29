package main

import (
	"io/ioutil"
	"path/filepath"

	config "github.com/carloszimm/stack-mining/configs"
	csvUtils "github.com/carloszimm/stack-mining/internal/csv"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

func main() {
	distPosts := make(map[string][]*types.Post)

	dirs, err := ioutil.ReadDir(config.DATA_EXPLORER_PATH)
	util.CheckError(err)

	for _, f := range dirs {
		if f.IsDir() && f.Name() != "consolidated sources" {
			files, err := ioutil.ReadDir(filepath.Join(config.DATA_EXPLORER_PATH, f.Name()))
			util.CheckError(err)

			distPosts[f.Name()] = csvUtils.ReadPostsCSVs(filepath.Join(config.DATA_EXPLORER_PATH, f.Name()), files)
		}
	}

	var allPosts []*types.Post
	for dist, posts := range distPosts {
		go csvUtils.WritePostsCSV(filepath.Join(config.CONSOLIDATED_SOURCES_PATH, dist+".csv"), posts)
		allPosts = append(allPosts, posts...)
	}
	csvUtils.WritePostsCSV(
		filepath.Join(config.CONSOLIDATED_SOURCES_PATH, "all.csv"),
		csvUtils.SortPosts(allPosts))
}

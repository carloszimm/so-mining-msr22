package csvUtils

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/carloszimm/stack-mining/types"
	"github.com/carloszimm/stack-mining/util"
	"github.com/gocarina/gocsv"
)

func ReadPostsCSVs(filesPath string, files []fs.FileInfo) []*types.Post {
	c := make(chan []*types.Post)

	for _, f := range files {
		go readPostsCSV(filepath.Join(filesPath, f.Name()), c)
	}

	var resultPosts [][]*types.Post

	for i := 0; i < len(files); i++ {
		posts := <-c
		resultPosts = append(resultPosts, posts)
	}

	return sortPosts(removeDuplicates(resultPosts))
}

func sortPosts(posts []*types.Post) []*types.Post {
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].Id < posts[j].Id
	})
	return posts
}

func removeDuplicates(postsArray [][]*types.Post) []*types.Post {
	resultSet := make(map[int]*types.Post)

	for _, posts := range postsArray {
		for _, post := range posts {
			if _, ok := resultSet[post.Id]; !ok {
				resultSet[post.Id] = post
			}
		}
	}

	result := make([]*types.Post, 0, len(resultSet))
	for _, post := range resultSet {
		result = append(result, post)
	}

	return result
}

func readPostsCSV(path string, c chan []*types.Post) {
	postsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	util.CheckError(err)
	defer postsFile.Close()

	posts := []*types.Post{}

	if err := gocsv.UnmarshalFile(postsFile, &posts); err != nil {
		panic(err)
	}

	c <- posts
}

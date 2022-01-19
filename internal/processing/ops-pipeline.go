package processing

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/dlclark/regexp2"
	"github.com/iancoleman/orderedmap"
)

// comment pattern acquired from:
// https://stackoverflow.com/questions/36725194/golang-regex-replace-excluding-quoted-strings

var (
	commentsReg = regexp2.MustCompile(
		`((?:(?:^[ \t]*)?(?:/\*[^*]*\*+(?:[^/*][^*]*\*+)*/(?:[ \t]*\r?\n(?=[ \t]*(?:\r?\n|/\*|//)))?|//(?:[^\\]|\\(?:\r?\n)?)*?(?:\r?\n(?=[ \t]*(?:\r?\n|/\*|//))|(?=\r?\n))))+)|("[^"\\]*(?:\\[\S\s][^"\\]*)*"|'[^'\\]*(?:\\[\S\s][^'\\]*)*'|(?:\r?\n|[\S\s])[^/"'\\\s]*)`, 0)
	stringsReg = regexp2.MustCompile(
		`(["'`+"`"+`])(?:(?=(\\?))\2.)*?\1`, 0)
)

func SetupOpsPipeline(posts []*types.Post, operators types.Operators) <-chan *orderedmap.OrderedMap {
	inOps, outOps := operators.CreateWorkerOps()

	out := createMsgPosts(posts)
	out = retriveTag(out)
	out = removeComments(out)
	out = removeStrings(out)
	dispatchToOpsCounters(out, inOps)
	return gatherResults(outOps)
}

func createMsgPosts(posts []*types.Post) <-chan types.PostMsg {
	out := make(chan types.PostMsg)
	go func() {
		for _, val := range posts {
			// transports only the posts' ids and bodies
			out <- types.PostMsg{PostId: val.Id, Body: val.Body}
		}
		close(out)
	}()
	return out
}

func retriveTag(in <-chan types.PostMsg) <-chan types.PostMsg {
	out := make(chan types.PostMsg)
	go func() {
		for postMsg := range in {
			// load the HTML document
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(postMsg.Body))
			util.CheckError(err)

			// retrieve code tags
			postMsg.Body = doc.Find("code").Text()

			out <- postMsg
		}
		close(out)
	}()
	return out
}

func removeComments(in <-chan types.PostMsg) <-chan types.PostMsg {
	out := make(chan types.PostMsg, 10)
	go func() {
		for postMsg := range in {
			// replaces comments by space
			postMsg.Body, _ = commentsReg.Replace(postMsg.Body, "$2 ", -1, -1)
			out <- postMsg
		}
		close(out)
	}()
	return out
}

func removeStrings(in <-chan types.PostMsg) <-chan types.PostMsg {
	out := make(chan types.PostMsg)
	go func() {
		for postMsg := range in {
			// replaces strings by an empty one
			postMsg.Body, _ = stringsReg.Replace(postMsg.Body, "", -1, -1)

			out <- postMsg
		}
		close(out)
	}()
	return out
}

// broadcast to operator counters(workers)
func dispatchToOpsCounters(in <-chan types.PostMsg, inOps []chan types.PostMsg) {
	go func() {
		for postMsg := range in {
			for _, inOp := range inOps {
				inOp <- postMsg
			}
		}
		// closes all channels when done
		types.CloseAllInOps(inOps)
	}()
}

func gatherResults(outOps chan types.CountMsg) <-chan *orderedmap.OrderedMap {
	out := make(chan *orderedmap.OrderedMap)
	result := orderedmap.New()
	go func() {
		for msg := range outOps {
			postId := fmt.Sprint(msg.PostId)
			if _, ok := result.Get(postId); !ok {
				// if doesn't exists yet, creat a new map
				result.Set(postId, orderedmap.New())
			}
			mI, _ := result.Get(postId)
			m := mI.(*orderedmap.OrderedMap)
			m.Set(msg.OperatorCount.Operator, msg.OperatorCount.Total)
		}
		// sort result by post Id
		result.SortKeys(sort.Strings)
		// sort results by operators' names
		keys := result.Keys()
		for _, k := range keys {
			valI, _ := result.Get(k)
			val := valI.(*orderedmap.OrderedMap)
			val.SortKeys(sort.Strings)
		}
		out <- result
		close(out)
	}()
	return out
}

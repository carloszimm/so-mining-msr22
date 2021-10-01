package processing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/dlclark/regexp2"
)

var (
	commentsReg = regexp.MustCompile(
		`(?:/\*[^*]*\*+(?:[^/*][^*]*\*+)*/|//[^\n]*(?:\n|$))|("[^"\\]*(?:\\[\S\s][^"\\]*)*"|'[^'\\]*(?:\\[\S\s][^'\\]*)*'|[\S\s][^/"'\\]*)`)
	stringsReg = regexp2.MustCompile(
		`(["'`+"`"+`])(?:(?=(\\?))\2.)*?\1`, 0)
)

func SetupOpsPipeline(posts []*types.Post, operators types.Operators) <-chan map[int][]types.OperatorCount {
	inOps, outOps := operators.CreateWorkerOps()

	out := createMsgPosts(posts)
	out = retriveTag(out)
	out = removeComments(out)
	out = removeStrings(out)
	dispatchToOpsCounters(out, inOps)
	return gatherResults(outOps, len(posts)*len(operators.GetOperators()))
}

func createMsgPosts(posts []*types.Post) <-chan types.PostMsg {
	out := make(chan types.PostMsg)
	go func() {
		for _, val := range posts {
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

			// retrieve code, pre, blockquotes tags
			postMsg.Body = doc.Find("code, pre, blockquotes").Text()

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
			postMsg.Body = commentsReg.ReplaceAllString(postMsg.Body, "$1")
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
			results := util.Regexp2FindAllString(stringsReg, postMsg.Body)

			for _, result := range results {
				postMsg.Body = strings.ReplaceAll(postMsg.Body, result[0].String(), "")
			}

			out <- postMsg
		}
		close(out)
	}()
	return out
}

func dispatchToOpsCounters(in <-chan types.PostMsg, inOps []chan types.PostMsg) {
	go func() {
		for postMsg := range in {
			for _, inOp := range inOps {
				inOp <- postMsg
			}
		}
	}()
}

func gatherResults(outOps chan types.CountMsg, totalMsgs int) <-chan map[int][]types.OperatorCount {
	out := make(chan map[int][]types.OperatorCount)
	result := make(map[int][]types.OperatorCount)
	go func() {
		for i := 0; i < totalMsgs; i++ {
			msg := <-outOps
			result[msg.PostId] = append(result[msg.PostId], msg.OperatorCount)
		}
		for _, val := range result {
			types.SortOperatorsCount(val)
		}
		out <- result
		close(out)
	}()
	return out
}

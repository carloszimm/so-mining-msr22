package types

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sort"
	"sync"

	config "github.com/carloszimm/stack-mining/configs"
	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/dlclark/regexp2"
)

// type to store the only fields of interest (PostId and Body)
type PostMsg struct {
	PostId int
	Body   string
}

// type to store the post id and some operator counting
type CountMsg struct {
	PostId int
	OperatorCount
}

type OperatorCount struct {
	Operator string
	Total    int
}

// type to hold operators' statistics (used only by the end of the process)
type OperatorStats struct {
	Sum    int
	Mean   float64
	StdDev float64
	Min    int
	Max    int
	Median int
}

// type used to hold operators' names and generate workers per request
type Operators struct {
	Dist          string
	operatorsList []string
}

func (ops *Operators) GetOperators() []string {
	return ops.operatorsList
}

func (ops *Operators) CreateWorkerOps() ([]chan PostMsg, chan CountMsg) {
	inChannels := make([]chan PostMsg, len(ops.operatorsList))
	outChannels := make([]chan CountMsg, len(ops.operatorsList))

	// ranges through the operators list and creates a channel/worker for each operator
	// stores output channels for futher merging and creation of a single channel
	for i, op := range ops.operatorsList {
		inChannel := make(chan PostMsg, 20)
		inChannels[i] = inChannel
		outChannels[i] = createOpWorker(inChannel, op)
	}

	return inChannels, mergeCountMsgs(outChannels...)
}

func createOpWorker(in <-chan PostMsg, opName string) chan CountMsg {
	out := make(chan CountMsg)
	counterFn := createCounter(opName)
	go func() {
		for postMsg := range in {
			out <- CountMsg{postMsg.PostId, OperatorCount{opName, counterFn(postMsg.Body)}}
		}
		close(out)
	}()
	return out
}

// based on https://go.dev/blog/pipelines
func mergeCountMsgs(cs ...chan CountMsg) chan CountMsg {
	var wg sync.WaitGroup
	out := make(chan CountMsg)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c chan CountMsg) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Dist is not used currently
func CreateOperators(path string, dist string) Operators {
	data, err := ioutil.ReadFile(filepath.Join(config.OPERATORS_PATH, path))
	util.CheckError(err)

	ops := Operators{}

	err = json.Unmarshal(data, &ops.operatorsList)
	util.CheckError(err)

	return ops
}

// (?<!\w) - negative look-behind to make sure the operator name isn't preceded by any character beside its own name
// \s* - followed by zero or more spaces
// { - for swift closures
func createCounter(opName string) func(string) int {
	re := regexp2.MustCompile(opName+`|\.?(?<!\w)`+opName+`\s*(\(|{)`, 0)
	return func(s string) int {
		return len(util.Regexp2FindAllString(re, s))
	}
}

// sort operator count by operators' names
func SortOperatorsCount(opCount []OperatorCount) {
	sort.SliceStable(opCount, func(i, j int) bool {
		return opCount[i].Operator < opCount[j].Operator
	})
}

func CloseAllInOps(inOps []chan PostMsg) {
	for _, in := range inOps {
		close(in)
	}
}

// returns a map containing operator name as key and an array of total counts as value
func AggregateByOperator(result interface{}) map[string][]int {
	opsCount := make(map[string][]int)
	switch r := result.(type) {
	case map[int][]OperatorCount:
		for _, val := range r {
			for _, opCount := range val {
				opsCount[opCount.Operator] = append(opsCount[opCount.Operator], opCount.Total)
			}
		}
	case map[string][]OperatorCount:
		for _, val := range r {
			for _, opCount := range val {
				opsCount[opCount.Operator] = append(opsCount[opCount.Operator], opCount.Total)
			}
		}
	}
	return opsCount
}

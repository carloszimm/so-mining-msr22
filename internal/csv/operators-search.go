package csvUtil

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

func WriteOpsSearchResult(path string, opsStats map[string]types.OperatorStats) {
	f, err := os.Create(path + ".csv")
	util.CheckError(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	p := sortMapByValue(opsStats)

	head := false
	for _, pair := range p {
		if !head {
			fields := []string{"Operator"}
			val := reflect.ValueOf(pair.Value)
			for i := 0; i < val.NumField(); i++ {
				fields = append(fields, val.Type().Field(i).Name)
			}
			err = w.Write(fields)
			util.CheckError(err)
			head = true
		}
		valFields := []string{pair.Key}
		val := reflect.ValueOf(pair.Value)
		for i := 0; i < val.NumField(); i++ {
			valFields = append(valFields, fmt.Sprintf("%v", val.Field(i).Interface()))
		}
		err = w.Write(valFields)
		util.CheckError(err)
	}

}

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value types.OperatorStats
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value.Sum > p[j].Value.Sum }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]types.OperatorStats) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

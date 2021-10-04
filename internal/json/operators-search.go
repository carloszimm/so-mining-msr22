package json

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/carloszimm/stack-mining/internal/util"
)

func WriteOpsSearchResult(wg *sync.WaitGroup, path string, opsStats map[string]types.OperatorStats) {
	defer wg.Done()

	j, err := json.Marshal(opsStats)
	util.CheckError(err)

	err = os.WriteFile(path+".json", j, 0644)
	util.CheckError(err)
}

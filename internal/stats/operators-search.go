package stats

import (
	"github.com/carloszimm/stack-mining/internal/types"
	"github.com/montanaflynn/stats"
)

func GenerateOpsStats(opsCount map[string][]int) map[string]types.OperatorStats {
	statsResult := make(map[string]types.OperatorStats)
	for operator, opCount := range opsCount {
		data := stats.LoadRawData(opCount)
		sum, _ := stats.Sum(data)
		mean, _ := stats.Mean(data)
		stdDev, _ := stats.StandardDeviation(data)
		min, _ := stats.Min(data)
		max, _ := stats.Max(data)
		median, _ := stats.Median(data)

		statsResult[operator] = types.OperatorStats{Sum: int(sum), Mean: mean,
			StdDev: stdDev, Min: int(min), Max: int(max), Median: int(median)}
	}

	return statsResult
}

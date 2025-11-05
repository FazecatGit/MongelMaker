package strategy

import "github.com/fazecat/mongelmaker/Internal/utils"

type VolumeStats struct {
	Timestamp   string
	TotalVolume int64
	MeanVolume  float64
	StdDev      float64
	Zscore      float64
	IsAnomalous bool
}

type WhaleEvent struct {
	Timestamp   string
	Symbol      string
	Direction   string // buy or sell
	Volume      int64
	ZScore      float64
	ClosePrice  float64
	PriceChange float64
	Conviction  string
}

func CalculateVolumeStats(volumes []int64) (mean float64, stdDev float64) {
	floatVolumes := make([]float64, len(volumes))
	for i, v := range volumes {
		floatVolumes[i] = float64(v)
	}

	mean = utils.Average(floatVolumes)
	stdDev = utils.StandardDeviation(floatVolumes)

	return mean, stdDev
}

func CalculateZScore(currentVolume int64, meanVolume float64, stdDev float64) float64 {
	if stdDev == 0 {
		return 0
	}

	zScore := (float64(currentVolume) - meanVolume) / stdDev
	return zScore
}

package utils

import "math"

func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}
func Abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func Max(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
func Min(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func CalculateAvgVolume(volumes []int64, period int) float64 {
	if len(volumes) < period {
		period = len(volumes)
	}
	sum := 0.0
	for i := len(volumes) - period; i < len(volumes); i++ {
		sum += float64(volumes[i])
	}
	return sum / float64(period)
}

func StandardDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := Average(values)
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance = variance / float64(len(values))
	return math.Sqrt(variance)
}

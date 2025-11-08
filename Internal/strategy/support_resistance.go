package strategy

import "github.com/fazecat/mongelmaker/Internal/types"

type PriceLevel struct {
	Price      float64
	BouncCount int
	Strength   float64
}

// finds the lowest support level
func FindSupport(bars []types.Bar) float64 {
	if len(bars) < 3 {
		return 0
	}
	lowestLow := bars[0].Low
	for _, bar := range bars {
		if bar.Low < lowestLow {
			lowestLow = bar.Low
		}
	}
	return lowestLow
}

// finds the highest resistance level
func FindResistance(bars []types.Bar) float64 {
	if len(bars) < 3 {
		return 0
	}
	highestHigh := bars[0].High
	for _, bar := range bars {
		if bar.High > highestHigh {
			highestHigh = bar.High
		}
	}
	return highestHigh
}

// finds multiple support levels (local minima/valleys)
func GetSupportLevels(bars []types.Bar) []PriceLevel {
	levels := []PriceLevel{}

	for i := 1; i < len(bars)-1; i++ {
		if bars[i].Low < bars[i-1].Low && bars[i].Low < bars[i+1].Low {
			levels = append(levels, PriceLevel{
				Price:      bars[i].Low,
				BouncCount: 1,
			})
		}
	}
	return levels
}

// finds multiple resistance levels (local maxima/peaks)
func GetResistanceLevels(bars []types.Bar) []PriceLevel {
	levels := []PriceLevel{}

	for i := 1; i < len(bars)-1; i++ {
		if bars[i].High > bars[i-1].High && bars[i].High > bars[i+1].High {
			levels = append(levels, PriceLevel{
				Price:      bars[i].High,
				BouncCount: 1,
			})
		}
	}
	return levels
}

// checks if price is near support
func IsAtSupport(currentPrice float64, support float64) bool {
	tolerance := support * 0.01
	return currentPrice >= support-tolerance && currentPrice <= support+tolerance
}

// checks if price is near resistance
func IsAtResistance(currentPrice float64, resistance float64) bool {
	tolerance := resistance * 0.01
	return currentPrice >= resistance-tolerance && currentPrice <= resistance+tolerance
}

// calculates % distance to support
func DistanceToSupport(currentPrice float64, support float64) float64 {
	if support == 0 {
		return 0
	}
	return ((currentPrice - support) / support) * 100
}

// calculates % distance to resistance
func DistanceToResistance(currentPrice float64, resistance float64) float64 {
	if resistance == 0 {
		return 0
	}
	return ((resistance - currentPrice) / resistance) * 100
}

//  calculates classic pivot point
func FindPivotPoint(bars []types.Bar) float64 {
	if len(bars) == 0 {
		return 0
	}

	latestBar := bars[0]
	return (latestBar.High + latestBar.Low + latestBar.Close) / 3
}

// checks if price breaks above resistance with buffer
func IsBreakoutAboveResistance(currentPrice float64, resistance float64) bool {
	if resistance == 0 {
		return false
	}
	return currentPrice > resistance*1.005 // 0.5% above resistance = breakout
}

// checks if price breaks below support with buffer
func IsBreakoutBelowSupport(currentPrice float64, support float64) bool {
	if support == 0 {
		return false
	}
	return currentPrice < support*0.995 // 0.5% below support = breakdown
}

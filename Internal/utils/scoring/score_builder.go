package scoring

import (
	"github.com/fazecat/mongelmaker/Internal/types"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

type ScoringInput struct {
	CurrentPrice float64
	VWAPPrice    float64
	ATRValue     float64
	RSIValue     float64
	WhaleCount   float64
	PriceDrop    float64
	ATRCategory  string
}

func BuildScoringInput(bars []types.Bar, vwapPrice float64, rsiValue float64, whaleCount int, atrValue float64, atrCategory string) (ScoringInput, error) {
	if len(bars) < 2 {
		return ScoringInput{}, nil
	}

	currentBar := bars[len(bars)-1]
	currentPrice := currentBar.Close
	openPrice := currentBar.Open

	priceDrop := 0.0
	if openPrice != 0 {
		priceDrop = ((openPrice - currentPrice) / openPrice) * 100
	}

	return ScoringInput{
		CurrentPrice: currentPrice,
		VWAPPrice:    vwapPrice,
		ATRValue:     atrValue,
		RSIValue:     rsiValue,
		WhaleCount:   float64(whaleCount),
		PriceDrop:    priceDrop,
		ATRCategory:  atrCategory,
	}, nil
}

func CalculateATRFromBars(bars []types.Bar) float64 {
	if len(bars) < 2 {
		return 0
	}

	trueRanges := make([]float64, len(bars))
	for i := 1; i < len(bars); i++ {
		high := bars[i].High
		low := bars[i].Low
		prevClose := bars[i-1].Close
		tr := utils.Max(high-low, utils.Abs(high-prevClose), utils.Abs(low-prevClose))
		trueRanges[i] = tr
	}

	period := 14
	if len(trueRanges) < period {
		period = len(trueRanges) - 1
	}

	return utils.Average(trueRanges[len(trueRanges)-period:])
}

func CategorizeATRValue(currentATR float64, bars []types.Bar) string {
	if len(bars) < 15 {
		return "NORMAL"
	}

	atrValues := make([]float64, 0)
	for i := len(bars) - 14; i < len(bars); i++ {
		subset := bars[:i+1]
		atr := CalculateATRFromBars(subset)
		if atr > 0 {
			atrValues = append(atrValues, atr)
		}
	}

	if len(atrValues) == 0 {
		return "NORMAL"
	}

	avgATR := utils.Average(atrValues)

	if currentATR < avgATR*0.5 {
		return "LOW"
	} else if currentATR > avgATR*1.5 {
		return "HIGH"
	}
	return "NORMAL"
}

func ScoreCategory(score float64) string {
	if score >= 8.0 {
		return "🟢 Excellent"
	}
	if score >= 6.0 {
		return "🟢 Good"
	}
	if score >= 4.0 {
		return "🟡 Fair"
	}
	if score >= 2.0 {
		return "🟠 Moderate"
	}
	return "🔴 Poor"
}

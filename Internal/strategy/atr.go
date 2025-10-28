package strategy

import (
	"fmt"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

func CalculateATR(atrBars []datafeed.ATRBar, period int) ([]float64, error) {
	if len(atrBars) < period+1 {
		return nil, fmt.Errorf("not enough data")
	}
	atrValues := make([]float64, len(atrBars))
	trueRanges := make([]float64, len(atrBars))

	// Calculate True Ranges
	for i := 1; i < len(atrBars); i++ {
		high := atrBars[i].High
		low := atrBars[i].Low
		prevClose := atrBars[i-1].Close

		trueRange := utils.Max(high-low, utils.Abs(high-prevClose), utils.Abs(low-prevClose))
		trueRanges[i] = trueRange
	}

	// Calculate ATR
	for i := period; i < len(atrBars); i++ {
		atrValues[i] = utils.Average(trueRanges[i-period+1 : i+1])
	}

	return atrValues, nil
}

func CalculateTrueRange(high, low, prevClose float64) float64 {
	return utils.Max(high-low, utils.Abs(high-prevClose), utils.Abs(low-prevClose))
}

func CalculateAndStoreATR(symbol string, period int, timeframe string, limit int) error {
	atrBars, err := datafeed.FetchATRPrices(symbol, limit, timeframe)
	if err != nil {
		return err
	}

	atrValues, err := CalculateATR(atrBars, period)
	if err != nil {
		return err
	}

	savedCount := 0
	for i, atrValue := range atrValues {
		if atrValue != 0 {
			err := datafeed.SaveATR(symbol, atrBars[i].Timestamp, atrValue)
			if err != nil {
				return err
			}
			savedCount++
		}
	}

	latestATR := 0.0
	for i := len(atrValues) - 1; i >= 0; i-- {
		if atrValues[i] != 0 {
			latestATR = atrValues[i]
			break
		}
	}

	fmt.Printf("✅ Saved %d ATR values for %s. Latest: %.2f\n", savedCount, symbol, latestATR)
	return nil
}

// Will do later
func DetermineATRSignal(atrValue float64, threshold float64) string {
	if atrValue > threshold {
		return "high volatility"
	}
	return "low volatility"
}

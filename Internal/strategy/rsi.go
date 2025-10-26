package strategy

import (
	"fmt"

	"github.com/fazecat/mongelmaker/Internal/utils"
)

func CalculateRSI(closes []float64, period int) ([]float64, error) {
	rsi := make([]float64, len(closes))
	if len(closes) < period {
		return nil, fmt.Errorf("not enough data to calculate RSI")
	}
	gains := make([]float64, len(closes))
	losses := make([]float64, len(closes))

	// Calculate gains and losses
	for i := 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	// Calculate average gains and losses
	averageGain := utils.Average(gains[len(gains)-period:])
	averageLoss := utils.Average(losses[len(losses)-period:])

	// Calculate RSI
	for i := period; i < len(closes); i++ {
		if averageLoss == 0 {
			rsi[i] = 100
		} else {
			rs := averageGain / averageLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi, nil
}

func FetchClosingPrices(symbol string, days int) ([]float64, error) {
	// Placeholder: Fetch closing prices from database or API
	return []float64{}, nil
}

func SaveRSI(symbol string, date string, rsiValue float64) error {
	// Placeholder: Save RSI value to database
	return nil
}
func DetermineRSISignal(rsiValue float64) string {
	// Placeholder: Determine buy/sell/hold based on RSI value
	return "hold"
}
func CalculateAndStoreRSI(symbol string) error {
	// Placeholder: Fetch closing prices
	return nil
}

package strategy

import (
	"fmt"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

func CalculateRSI(closes []float64, period int) ([]float64, error) {

	if len(closes) < period+1 {
		return nil, fmt.Errorf("not enough data")
	}
	rsi := make([]float64, len(closes))

	// Calculate gains and losses
	gains := make([]float64, len(closes))
	losses := make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	for i := period; i < len(closes); i++ {
		windowGains := gains[i-period+1 : i+1]
		windowLosses := losses[i-period+1 : i+1]

		// Calculate average window
		avgGain := utils.Average(windowGains)
		avgLoss := utils.Average(windowLosses)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi, nil
}

func DetermineRSISignal(rsiValue float64) string {
	if rsiValue < 30 {
		return "oversold"
	} else if rsiValue > 70 {
		return "overbought"
	}
	return "neutral"
}

func CalculateAndStoreRSI(symbol string, period int, timeframe string, limit int) error {
	pricePoints, err := datafeed.FetchPricePoints(symbol, limit, timeframe)
	if err != nil {
		return fmt.Errorf("failed to fetch price points: %w", err)
	}

	closes := make([]float64, len(pricePoints))
	for i, pp := range pricePoints {
		closes[i] = pp.Price
	}

	rsiValues, err := CalculateRSI(closes, period)
	if err != nil {
		return fmt.Errorf("failed to calculate RSI: %w", err)
	}

	for i := period; i < len(pricePoints); i++ {
		err = datafeed.SaveRSI(symbol,
			pricePoints[i].Timestamp,
			rsiValues[i])
		if err != nil {
			return fmt.Errorf("failed to save RSI for timestamp %s: %w",
				pricePoints[i].Timestamp.Format("2006-01-02 15:04:05"), err)
		}
	}

	latestRSI := rsiValues[len(rsiValues)-1]
	fmt.Printf("✅ Saved %d RSI values for %s. Latest: %.2f\n", len(pricePoints)-period, symbol, latestRSI)
	return nil
}

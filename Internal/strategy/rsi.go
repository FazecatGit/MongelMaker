package strategy

import (
	"fmt"
	"time"

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

func CalculateAndStoreRSI(symbol string, period int) error {
	closes, err := datafeed.FetchClosingPrices(symbol, period*3)
	if err != nil {
		return fmt.Errorf("failed to fetch prices: %w", err)
	}

	rsiValues, err := CalculateRSI(closes, period)
	if err != nil {
		return fmt.Errorf("failed to calculate RSI: %w", err)
	}
	defer datafeed.CloseDatabase()

	latesetRSI := rsiValues[len(rsiValues)-1]

	err = datafeed.SaveRSI(symbol, time.Now().Format("2006-01-02"), latesetRSI)
	if err != nil {
		return fmt.Errorf("failed to save RSI: %w", err)
	}

	fmt.Printf("RSI for %s on latest date: %.2f\n", symbol, latesetRSI)
	return nil
}

package datafeed

import (
	"fmt"
)

type Signal struct {
	symbol       string
	currentPrice float64
	SMA          float64
	Action       string // "BUY", "SELL", "HOLD"
}

func GetCurrentPrice(symbol string) (float64, error) {
	Quote, err := GetLastQuote(symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to grab last quote for %s: %w", symbol, err)
	}
	return Quote.Price, nil
}

func HistoricalData(symbol string, timeframe string, limit int) ([]Bar, error) {
	limitbars, err := GetAlpacaBars(symbol, timeframe, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data for %s: %w", symbol, err)
	}
	oneminutebars, err := GetAlpacaBars(symbol, "1Min", 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1Min bars for %s: %w", symbol, err)
	}
	fiveminutebars, err := GetAlpacaBars(symbol, "5Min", 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get 5Min bars for %s: %w", symbol, err)
	}
	onedaybars, err := GetAlpacaBars(symbol, "1Day", 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1Day bars for %s: %w", symbol, err)
	}
	return append(append(append(limitbars, oneminutebars...), fiveminutebars...), onedaybars...), nil
	//append inception wtf
}

func calculateSMA(bars []Bar) float64 {
	if len(bars) == 0 {
		return 0
	}
	sum := 0.0
	for _, bar := range bars {
		sum += bar.Close
	}
	return sum / float64(len(bars))
}

func GenerateSignal(symbol string) (*Signal, error) {
	// Fetch current price
	currentPrice, err := GetCurrentPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get current price for %s: %w", symbol, err)
	}

	// Fetch historical data
	historicalData, err := HistoricalData(symbol, "1Day", 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data for %s: %w", symbol, err)
	}

	// Calculate SMA
	sma := calculateSMA(historicalData)

	// Generate signal
	var action string
	if currentPrice > sma {
		action = "BUY"
	} else if currentPrice < sma {
		action = "SELL"
	} else {
		action = "HOLD"
	}

	return &Signal{
		symbol:       symbol,
		currentPrice: currentPrice,
		SMA:          sma,
		Action:       action,
	}, nil
}

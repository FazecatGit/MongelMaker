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
	return []Bar{limitbars, oneminutebars, fiveminutebars, onedaybars}, nil
}

func calculateSMA(bars []Bar) float64 {
	//placeholder for SMA calculation
}

func GenerateSignal(symbol string) (*Signal, error) {
	// Fetch current price
}

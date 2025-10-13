package datafeed

import (
	"fmt"
	"time"
)

func StoreBarsWithAnalytics(symbol string, timeframe string, bars []Bar) error {
	query := `
		INSERT INTO historical_bars (symbol, timeframe, timestamp, open_price, high_price, low_price, close_price, 
		volume, price_change, price_change_percent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (symbol, timeframe, timestamp) DO NOTHING;
	`
	for _, bar := range bars {
		timestamp, err := time.Parse(time.RFC3339, bar.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp %s: %w", bar.Timestamp, err)
		}

		priceChange := bar.Close - bar.Open
		priceChangePercent := (bar.Close - bar.Open) / bar.Open * 100

		_, err = DB.Exec(query, symbol, timeframe, timestamp, bar.Open, bar.High, bar.Low,
			bar.Close, bar.Volume, priceChange, priceChangePercent)
		if err != nil {
			return fmt.Errorf("failed to insert bar: %w", err)
		}
	}
	return nil
}

package datafeed

import (
	"context"
	"strconv"
	"time"

	database "github.com/fazecat/mongelmaker/Internal/database/sqlc"
)

type PricePoint struct {
	Price     float64
	Timestamp time.Time
}

func FetchClosingPrices(symbol string, days int, timeframe string) ([]float64, error) {
	params := database.GetClosingPricesParams{
		Symbol:    symbol,
		Timeframe: timeframe,
		Limit:     int32(days),
	}

	ctx := context.Background()
	rows, err := Queries.GetClosingPrices(ctx, params)
	if err != nil {
		return nil, err
	}

	var closingPrices []float64
	for _, row := range rows {
		price, err := strconv.ParseFloat(row.ClosePrice, 64)
		if err != nil {
			return nil, err
		}
		closingPrices = append(closingPrices, price)
	}

	return closingPrices, nil
}
func SaveRSI(symbol string, timestamp time.Time, rsiValue float64) error {
	params := database.SaveRSIParams{
		Symbol:               symbol,
		CalculationTimestamp: timestamp,
		RsiValue:             float32(rsiValue),
	}
	ctx := context.Background()
	err := Queries.SaveRSI(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func FetchPricePoints(symbol string, days int, timeframe string) ([]PricePoint, error) {
	params := database.GetClosingPricesParams{
		Symbol:    symbol,
		Timeframe: timeframe,
		Limit:     int32(days),
	}
	ctx := context.Background()
	rows, err := Queries.GetClosingPrices(ctx, params)
	if err != nil {
		return nil, err
	}

	var pricePoints []PricePoint
	for _, row := range rows {
		price, err := strconv.ParseFloat(row.ClosePrice, 64)
		if err != nil {
			return nil, err
		}
		pricePoints = append(pricePoints, PricePoint{
			Price:     price,
			Timestamp: row.Timestamp,
		})
	}

	return pricePoints, nil
}

func FetchRSIForDisplay(symbol string, limit int) (map[string]float64, error) {
	params := database.GetRSIForDateRangeParams{
		Symbol: symbol,
		Limit:  int32(limit),
	}
	ctx := context.Background()
	rows, err := Queries.GetRSIForDateRange(ctx, params)
	if err != nil {
		return nil, err
	}

	rsiMap := make(map[string]float64)
	for _, row := range rows {
		dateStr := row.CalculationTimestamp.Format("2006-01-02 15:04:05")
		rsiMap[dateStr] = float64(row.RsiValue)
	}
	return rsiMap, nil
}

func FetchRSIByTimestampRange(symbol string, startTime, endTime time.Time) (map[string]float64, error) {
	params := database.GetRSIByTimestampRangeParams{
		Symbol:                 symbol,
		CalculationTimestamp:   startTime,
		CalculationTimestamp_2: endTime,
	}
	ctx := context.Background()
	rows, err := Queries.GetRSIByTimestampRange(ctx, params)
	if err != nil {
		return nil, err
	}

	rsiMap := make(map[string]float64)
	for _, row := range rows {
		dateStr := row.CalculationTimestamp.Format("2006-01-02 15:04:05")
		rsiMap[dateStr] = float64(row.RsiValue)
	}
	return rsiMap, nil
}

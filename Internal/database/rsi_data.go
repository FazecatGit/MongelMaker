package datafeed

import (
	"context"
	"fmt"
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

func SaveRSI(symbol string, date string, rsiValue float64) error {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	params := database.SaveRSIParams{
		Symbol:          symbol,
		CalculationDate: parsedDate,
		RsiValue:        strconv.FormatFloat(rsiValue, 'f', 2, 64),
	}
	ctx := context.Background()
	err = Queries.SaveRSI(ctx, params)
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
		dateStr := row.CalculationDate.Format("2006-01-02")
		rsiVal, _ := strconv.ParseFloat(row.RsiValue, 64)
		rsiMap[dateStr] = rsiVal
	}
	return rsiMap, nil
}

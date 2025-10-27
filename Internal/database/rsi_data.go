package datafeed

import (
	"context"
	"fmt"
	"strconv"
	"time"

	database "github.com/fazecat/mongelmaker/Internal/database/sqlc"
)

func FetchClosingPrices(symbol string, days int) ([]float64, error) {
	params := database.GetClosingPricesParams{
		Symbol: symbol,
		Limit:  int32(days),
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

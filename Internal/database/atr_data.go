package datafeed

import (
	"context"
	"strconv"
	"time"

	database "github.com/fazecat/mongelmaker/Internal/database/sqlc"
)

type ATRPoint struct {
	ATR       float64
	Timestamp time.Time
}

type ATRBar struct {
	High      float64
	Low       float64
	Close     float64
	Timestamp time.Time
}

func FetchATRPrices(symbol string, limit int) ([]ATRBar, error) {
	params := database.GetATRPricesParams{
		Symbol: symbol,
		Limit:  int32(limit),
	}
	ctx := context.Background()
	rows, err := Queries.GetATRPrices(ctx, params)

	var atrBars []ATRBar
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		atrHigh, err := strconv.ParseFloat(row.HighPrice, 64)
		if err != nil {
			return nil, err
		}
		atrLow, err := strconv.ParseFloat(row.LowPrice, 64)
		if err != nil {
			return nil, err
		}
		atrClose, err := strconv.ParseFloat(row.ClosePrice, 64)
		if err != nil {
			return nil, err
		}

		atrBars = append(atrBars, ATRBar{
			High:      atrHigh,
			Low:       atrLow,
			Close:     atrClose,
			Timestamp: row.Timestamp,
		})
	}
	return atrBars, nil
}

func SaveATR(symbol string, date string, atrValue float64) error {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}

	atrValueStr := strconv.FormatFloat(atrValue, 'f', -1, 64)

	params := database.SaveATRParams{
		Symbol:          symbol,
		CalculationDate: parsedDate,
		AtrValue:        atrValueStr,
	}
	ctx := context.Background()
	err = Queries.SaveATR(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

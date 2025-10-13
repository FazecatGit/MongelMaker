package interactive

import (
	"fmt"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
)

func ShowTimeframeMenu() (string, error) {
	fmt.Println("Choose timeframe:")
	fmt.Println("1. 1 Minute")
	fmt.Println("2. 5 Minutes")
	fmt.Println("3. 1 Hour")
	fmt.Println("4. 1 Day")

	fmt.Print("Enter choice: ")
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number between 1 and 4.")
		return "", err
	}

	switch choice {
	case 1:
		return "1Min", nil
	case 2:
		return "5Min", nil
	case 3:
		return "1Hour", nil
	case 4:
		return "1Day", nil
	default:
		fmt.Println("Invalid choice.")
	}
	return "", fmt.Errorf("invalid choice")
}

func FetchMarketData(symbol string, timeframe string, limit int) ([]datafeed.Bar, error) {
	if timeframe == "" {
		return nil, fmt.Errorf("timeframe cannot be empty")
	}
	bars, err := datafeed.GetAlpacaBars(symbol, timeframe, limit)
	if err != nil {
		return nil, err
	}
	return bars, nil
}

func DisplayBasicData(bars []datafeed.Bar, symbol string, timeframe string) {
	fmt.Printf("\n📊 Basic Data for %s (%s)\n", symbol, timeframe)
	fmt.Println("Timestamp           | Close Price | Volume")
	fmt.Println("--------------------|-------------|----------")

	for _, bar := range bars {
		fmt.Printf("%-20s | %11.2f | %8d\n", bar.Timestamp, bar.Close, bar.Volume)

	}
}

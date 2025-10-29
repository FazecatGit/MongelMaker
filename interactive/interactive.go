package interactive

import (
	"fmt"
	"time"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/strategy"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

func ShowTimeframeMenu() (string, error) {
	fmt.Println("Choose timeframe:")
	fmt.Println("1.  1 Minute")
	fmt.Println("2.  3 Minutes")
	fmt.Println("3.  5 Minutes")
	fmt.Println("4.  10 Minutes")
	fmt.Println("5.  30 Minutes")
	fmt.Println("6.  1 Hour")
	fmt.Println("7.  2 Hours")
	fmt.Println("8.  4 Hours")
	fmt.Println("9.  1 Day")
	fmt.Println("10. 1 Week")
	fmt.Println("11. 1 Month")
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
		return "3Min", nil
	case 3:
		return "5Min", nil
	case 4:
		return "10Min", nil
	case 5:
		return "30Min", nil
	case 6:
		return "1Hour", nil
	case 7:
		return "2Hour", nil
	case 8:
		return "4Hour", nil
	case 9:
		return "1Day", nil
	case 10:
		return "1Week", nil
	case 11:
		return "1Month", nil
	default:
		fmt.Println("Invalid choice.")
		return "", fmt.Errorf("invalid choice")
	}
}

func FetchMarketData(symbol string, timeframe string, limit int, startDate string) ([]datafeed.Bar, error) {
	if timeframe == "" {
		return nil, fmt.Errorf("timeframe cannot be empty")
	}

	// fetch enough bars for RSI/ATR calculations
	if limit < 14 {
		limit = 14
	}

	bars, err := datafeed.GetAlpacaBars(symbol, timeframe, limit, startDate)
	if err != nil {
		return nil, err
	}

	if len(bars) < 14 {
		return nil, fmt.Errorf("not enough data to calculate RSI/ATR: fetched %d bars, need at least 14", len(bars))
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

func DisplayAdvancedData(bars []datafeed.Bar, symbol string, timeframe string) {
	fmt.Printf("\n📊 Advanced Data for %s (%s)\n", symbol, timeframe)
	fmt.Println("Timestamp           | Open Price | High Price | Low Price | Close Price | Volume")
	fmt.Println("--------------------|------------|------------|-----------|-------------|----------")

	for _, bar := range bars {
		fmt.Printf("%-20s | %11.2f | %11.2f | %9.2f | %11.2f | %8d\n",
			bar.Timestamp, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
	}
}

func DisplayAnalyticsData(bars []datafeed.Bar, symbol string, timeframe string, tz *time.Location) {
	fmt.Printf("\n📈 Analytics Data for %s (%s) - Timezone: %s\n", symbol, timeframe, tz.String())

	var startTime, endTime time.Time
	if len(bars) > 0 {
		firstBar, err := time.Parse(time.RFC3339, bars[0].Timestamp)
		if err == nil {
			startTime = firstBar
		}
		lastBar, err := time.Parse(time.RFC3339, bars[len(bars)-1].Timestamp)
		if err == nil {
			endTime = lastBar
		}
	}

	// Fetch RSI and ATR data
	var rsiMap map[string]float64
	var atrMap map[string]float64
	var err error

	if !startTime.IsZero() && !endTime.IsZero() {
		rsiMap, err = datafeed.FetchRSIByTimestampRange(symbol, startTime, endTime)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch RSI data: %v\n", err)
			rsiMap = make(map[string]float64)
		}

		atrMap, err = datafeed.FetchATRByTimestampRange(symbol, startTime, endTime)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch ATR data: %v\n", err)
			atrMap = make(map[string]float64)
		}
	} else {
		// Fallback
		fetchLimit := len(bars) * 10
		rsiMap, err = datafeed.FetchRSIForDisplay(symbol, fetchLimit)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch RSI data: %v\n", err)
			rsiMap = make(map[string]float64)
		}

		atrMap, err = datafeed.FetchATRForDisplay(symbol, fetchLimit)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch ATR data: %v\n", err)
			atrMap = make(map[string]float64)
		}
	}

	fmt.Println("Timestamp           | Close Price | Price Chg | Chg %  | Volume   | RSI    | ATR    | B/U Ratio | B/L Ratio | Analysis                  | Signals             ")
	fmt.Println("--------------------|-------------|-----------|--------|----------|--------|--------|-----------|-----------|--------------------------|---------------------")

	for _, bar := range bars {
		priceChange := bar.Close - bar.Open
		priceChangePercent := (bar.Close - bar.Open) / bar.Open * 100

		t, err := time.Parse(time.RFC3339, bar.Timestamp)
		if err != nil {
			fmt.Printf("⚠️  Could not parse timestamp: %v\n", err)
		}

		var timestampStr string
		var displayTimestamp string

		if err == nil {
			localTime := t.In(tz)
			displayTimestamp = localTime.Format("2006-01-02 15:04:05")
			timestampStr = t.Format("2006-01-02 15:04:05")
		} else {
			timestampStr = bar.Timestamp
			displayTimestamp = bar.Timestamp
		}

		// Get RSI and ATR for current timestamp
		rsiVal, hasRSI := rsiMap[timestampStr]
		atrVal, hasATR := atrMap[timestampStr]

		rsiStr := "  -   "
		if hasRSI {
			rsiStr = fmt.Sprintf("%6.2f", rsiVal)
		}

		atrStr := "  -   "
		if hasATR {
			atrStr = fmt.Sprintf("%6.2f", atrVal)
		}

		// Calculate Body-to-Wick ratios and analysis
		candle := utils.Candlestick{
			Open:  bar.Open,
			Close: bar.Close,
			High:  bar.High,
			Low:   bar.Low,
		}
		metrics, results := utils.AnalyzeCandlestick(candle)
		bodyToUpperStr := fmt.Sprintf("%9.2f", metrics["BodyToUpper"])
		bodyToLowerStr := fmt.Sprintf("%9.2f", metrics["BodyToLower"])
		analysisStr := results["Analysis"]

		// visualize signals
		signalStr := ""
		if hasRSI {
			rsiSignal := strategy.DetermineRSISignal(rsiVal)
			switch rsiSignal {
			case "overbought":
				signalStr += "📈 Overbought"
			case "oversold":
				signalStr += "📉 Oversold"
			case "neutral":
				signalStr += "➡️  Neutral"
			}
		}

		if hasATR {
			// Calculate ATR threshold dynamically
			atrThreshold := bar.Close * 0.01
			atrSignal := strategy.DetermineATRSignal(atrVal, atrThreshold)

			if signalStr != "" {
				signalStr += " | "
			}

			switch atrSignal {
			case "high volatility":
				signalStr += "⚡ High Vol"
			case "low volatility":
				signalStr += "❄️  Low Vol"
			}
		}

		if signalStr == "" {
			signalStr = "-"
		}

		fmt.Printf("%-20s | %11.2f | %9.2f | %6.2f | %8d | %6s | %6s | %9s | %9s | %-25s | %-20s\n",
			displayTimestamp, bar.Close, priceChange, priceChangePercent, bar.Volume, rsiStr, atrStr, bodyToUpperStr, bodyToLowerStr, analysisStr, signalStr)
	}
}

func ShowDisplayMenu() (string, error) {
	fmt.Println("\nChoose display format:")
	fmt.Println("1. Basic Data")
	fmt.Println("2. Full OHLC")
	fmt.Println("3. Analytics")
	fmt.Println("4. All Data")

	fmt.Print("Enter choice: ")
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number between 1 and 4.")
		return "", err
	}

	switch choice {
	case 1:
		return "basic", nil
	case 2:
		return "full", nil
	case 3:
		return "analytics", nil
	case 4:
		return "all", nil
	default:
		fmt.Println("Invalid choice.")
	}
	return "", fmt.Errorf("invalid choice")
}

func ShowTimezoneMenu() (*time.Location, error) {
	nyLocation, _ := time.LoadLocation("America/New_York")
	chicagoLocation, _ := time.LoadLocation("America/Chicago")
	laLocation, _ := time.LoadLocation("America/Los_Angeles")
	londonLocation, _ := time.LoadLocation("Europe/London")
	tokyoLocation, _ := time.LoadLocation("Asia/Tokyo")
	hkLocation, _ := time.LoadLocation("Asia/Hong_Kong")

	timezones := map[int]struct {
		name     string
		location *time.Location
	}{
		1: {"UTC", time.UTC},
		2: {"America/New_York (NYSE/NASDAQ)", nyLocation},
		3: {"America/Chicago (CME)", chicagoLocation},
		4: {"America/Los_Angeles (PST)", laLocation},
		5: {"Europe/London (LSE)", londonLocation},
		6: {"Asia/Tokyo (TSE)", tokyoLocation},
		7: {"Asia/Hong_Kong (HKEX)", hkLocation},
		8: {"Local (System Time)", time.Local},
	}

	fmt.Println("\nChoose timezone:")
	for i := 1; i <= len(timezones); i++ {
		fmt.Printf("%d. %s\n", i, timezones[i].name)
	}

	fmt.Print("Enter choice: ")
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid input. Please enter a valid number.")
		return nil, err
	}

	if tz, exists := timezones[choice]; exists {
		return tz.location, nil
	} else {
		fmt.Println("Invalid choice. Defaulting to UTC.")
		return time.UTC, nil
	}
}

package datafeed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fazecat/mongelmaker/Internal/types"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

type Bar = types.Bar

func GetAlpacaBars(symbol string, timeframe string, limit int, startDate string) ([]Bar, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_API_SECRET")

	if startDate == "" {
		// Compute a start time far enough in the past to return `limit` bars
		now := time.Now().UTC()

		timeframeToDur := func(tf string) time.Duration {
			switch tf {
			case "1Min":
				return time.Minute
			case "3Min":
				return 3 * time.Minute
			case "5Min":
				return 5 * time.Minute
			case "10Min":
				return 10 * time.Minute
			case "30Min":
				return 30 * time.Minute
			case "1Hour":
				return time.Hour
			case "2Hour":
				return 2 * time.Hour
			case "4Hour":
				return 4 * time.Hour
			case "1Day":
				return 24 * time.Hour
			case "1Week":
				return 7 * 24 * time.Hour
			case "1Month":
				return 30 * 24 * time.Hour
			default:
				return 24 * time.Hour
			}
		}

		barDur := timeframeToDur(timeframe)
		// add a small safety buffer of 2 bars
		totalDur := barDur * time.Duration(limit+2)
		start := now.Add(-totalDur)
		startDate = start.Format(time.RFC3339)
	}

	apiURL := fmt.Sprintf(
		"https://data.alpaca.markets/v2/stocks/%s/bars?timeframe=%s&limit=%d&start=%s",
		symbol, timeframe, limit, startDate,
	)

	fmt.Printf("🔗 API Request: %s\n", apiURL)

	//retry logic
	var bars []Bar
	retryConfig := utils.DefaultRetryConfig()

	err := utils.RetryWithBackoff(func() error {
		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("APCA-API-KEY-ID", apiKey)
		req.Header.Set("APCA-API-SECRET-KEY", secretKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Printf("📡 API Response Status: %s\n", resp.Status)

		// Checking API errors
		if resp.StatusCode == 403 {
			fmt.Printf("⚠️  403 Forbidden - Your account may not have access to %s data\n", timeframe)
			bars = []Bar{} // Return empty slice instead of error
			return nil
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("API returned status %d", resp.StatusCode)
		}

		type Response struct {
			Bars []Bar `json:"bars"`
		}

		var r Response
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return err
		}

		bars = r.Bars
		return nil
	}, retryConfig)

	if err != nil {
		return nil, err
	}

	fmt.Printf("📊 Received %d bars\n", len(bars))

	// Reverse bars to latest-first (most recent data first)
	for i, j := 0, len(bars)-1; i < j; i, j = i+1, j-1 {
		bars[i], bars[j] = bars[j], bars[i]
	}

	return bars, nil
}

type LastQuote struct {
	Price float64 `json:"ap"`
}

func GetLastQuote(symbol string) (*LastQuote, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_API_SECRET")

	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/quotes/latest", url.PathEscape(symbol))

	var quote *LastQuote
	retryConfig := utils.DefaultRetryConfig()

	err := utils.RetryWithBackoff(func() error {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("APCA-API-KEY-ID", apiKey)
		req.Header.Set("APCA-API-SECRET-KEY", secretKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get last quote: %s", resp.Status)
		}

		type Response struct {
			Quote LastQuote `json:"quote"`
		}

		var r Response
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return err
		}

		quote = &r.Quote
		return nil
	}, retryConfig)

	return quote, err
}

func GetLastTrade(symbol string) (*Bar, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_API_SECRET")

	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/trades/latest", url.PathEscape(symbol))

	var trade *Bar
	retryConfig := utils.DefaultRetryConfig()

	err := utils.RetryWithBackoff(func() error {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("APCA-API-KEY-ID", apiKey)
		req.Header.Set("APCA-API-SECRET-KEY", secretKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get last trade: %s", resp.Status)
		}

		var r Bar
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return err
		}

		trade = &r
		return nil
	}, retryConfig)

	return trade, err
}

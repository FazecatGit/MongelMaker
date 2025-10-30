package datafeed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fazecat/mongelmaker/Internal/utils"
)

type Bar struct {
	Timestamp string  `json:"t"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    int64   `json:"v"`
}

func GetAlpacaBars(symbol string, timeframe string, limit int, startDate string) ([]Bar, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_API_SECRET")

	if startDate == "" {
		// Default to 90 days
		startDate = time.Now().UTC().AddDate(0, 0, -90).Format(time.RFC3339)
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

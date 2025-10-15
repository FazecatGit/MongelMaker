package datafeed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Bar struct {
	Timestamp string  `json:"t"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    int64   `json:"v"`
}

func GetAlpacaBars(symbol string, timeframe string, limit int) ([]Bar, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_SECRET_KEY")

	// Use simpler URL for now - some accounts may not have access to historical daily data
	apiURL := fmt.Sprintf(
		"https://data.alpaca.markets/v2/stocks/%s/bars?timeframe=%s&limit=%d",
		symbol, timeframe, limit,
	)

	fmt.Printf("🔗 API Request: %s\n", apiURL)

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("APCA-API-KEY-ID", apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("📡 API Response Status: %s\n", resp.Status)

	// Check for API errors
	if resp.StatusCode == 403 {
		fmt.Printf("⚠️  403 Forbidden - Your account may not have access to %s data\n", timeframe)
		return []Bar{}, nil // Return empty slice instead of error
	}

	type Response struct {
		Bars []Bar `json:"bars"`
	}

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	fmt.Printf("📊 Received %d bars\n", len(r.Bars))
	return r.Bars, nil
}

type LastQuote struct {
	Price float64 `json:"ap"`
}

func GetLastQuote(symbol string) (*LastQuote, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_SECRET_KEY")

	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/quotes/latest", url.PathEscape(symbol))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("APCA-API-KEY-ID", apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get last quote: %s", resp.Status)
	}

	type Response struct {
		Quote LastQuote `json:"quote"`
	}

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r.Quote, nil
}

func GetLastTrade(symbol string) (*Bar, error) {
	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_SECRET_KEY")

	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/trades/latest", url.PathEscape(symbol))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("APCA-API-KEY-ID", apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get last trade: %s", resp.Status)
	}

	var r Bar
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

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

	url := fmt.Sprintf(
		"https://data.alpaca.markets/v2/stocks/%s/bars?timeframe=%s&limit=%d",
		symbol, timeframe, limit,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("APCA-API-KEY-ID", apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Response struct {
		Bars []Bar `json:"bars"`
	}

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r.Bars, nil
}

func GetLastQuote(symbol string) (*struct {
	Price float64 `json:"ap"}`
}, error) {
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

	var quote struct {
		Price float64 `json:"ap"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/strategy"
	"github.com/fazecat/mongelmaker/interactive"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = datafeed.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer datafeed.CloseDatabase()

	// Test the retry logic
	// utils.TestRetryLogic()
	// "github.com/fazecat/mongelmaker/Internal/utils"

	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_API_SECRET")

	alpclient := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    apiKey,
		APISecret: secretKey,
		BaseURL:   "https://paper-api.alpaca.markets",
	})

	req, _ := http.NewRequest("GET", "https://paper-api.alpaca.markets/v2/account", nil)
	req.Header.Set("APCA-API-KEY-ID", apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", secretKey)

	account, err := alpclient.GetAccount()
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	timeframe, err := interactive.ShowTimeframeMenu()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get Stock Symbol from User
	var symbol string
	fmt.Print("\nEnter stock symbol to analyze (e.g., AAPL, TSLA, MSFT): ")
	fmt.Scan(&symbol)

	// Get number of bars to fetch/analyze from user
	var numBars int
	fmt.Print("\nHow many bars to fetch and analyze?: ")
	fmt.Scan(&numBars)

	// Calculate indicators and fetch data
	testIndicators(symbol, numBars, timeframe)

	bars, err := interactive.FetchMarketData(symbol, timeframe, numBars, "")
	if err != nil {
		fmt.Println("Error fetching market data:", err)
		return
	}

	displayChoice, err := interactive.ShowDisplayMenu()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	timezone, err := interactive.ShowTimezoneMenu()
	if err != nil {
		fmt.Println("Error selecting timezone, using UTC:", err)
		timezone = time.UTC
	}

	// Call the appropriate display function based on choice
	switch displayChoice {
	case "basic":
		interactive.DisplayBasicData(bars, symbol, timeframe)
	case "full":
		interactive.DisplayAdvancedData(bars, symbol, timeframe)
	case "analytics":
		interactive.DisplayAnalyticsData(bars, symbol, timeframe, timezone)
	case "all":
		interactive.DisplayBasicData(bars, symbol, timeframe)
		interactive.DisplayAdvancedData(bars, symbol, timeframe)
		interactive.DisplayAnalyticsData(bars, symbol, timeframe, timezone)
	}

	fmt.Printf("\n✅ Displayed %d bars for %s on %s timeframe\n", len(bars), symbol, timeframe)

	balanceChange := account.Equity.Sub(account.LastEquity)

	fmt.Println("Status:", resp.Status, balanceChange)
}

func testIndicators(symbol string, numBars int, timeframe string) {
	log.Println("🧪 Calculating RSI and ATR indicators...")

	bars, err := datafeed.GetAlpacaBars(symbol, timeframe, numBars, "")
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		return
	}
	log.Printf("✅ Fetched %d bars", len(bars))

	err = datafeed.StoreBarsWithAnalytics(symbol, timeframe, bars)
	if err != nil {
		log.Printf("Failed to store bars: %v", err)
		return
	}
	log.Println("✅ Stored bars in database")

	err = strategy.CalculateAndStoreRSI(symbol, 14, timeframe, numBars)
	if err != nil {
		log.Printf("Failed to calculate RSI: %v", err)
		return
	}
	log.Println("✅ RSI calculation and storage successful!")

	err = strategy.CalculateAndStoreATR(symbol, 14, timeframe, numBars)
	if err != nil {
		log.Printf("Failed to calculate ATR: %v", err)
		return
	}
	log.Println("✅ ATR calculation and storage successful!")
}

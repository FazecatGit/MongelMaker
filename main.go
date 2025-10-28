package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/strategy"
	"github.com/fazecat/mongelmaker/interactive"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
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

	// Get number of days from user
	var days int
	fmt.Print("How many days of history to analyze (recommended: 30-90): ")
	fmt.Scan(&days)

	// Get Stock Symbol from User
	var symbol string
	fmt.Print("\nEnter stock symbol to analyze (e.g., AAPL, TSLA, MSFT): ")
	fmt.Scan(&symbol)

	// Test the indicators of RSI and ATR
	testIndicators(symbol, days, timeframe)

	// Get number of bars to fetch from user
	var numBars int
	fmt.Print("How many bars to display? ")
	fmt.Scan(&numBars)

	bars, err := interactive.FetchMarketData(symbol, timeframe, numBars)
	if err != nil {
		fmt.Println("Error fetching market data:", err)
		return
	}

	// Get user's display preference
	displayChoice, err := interactive.ShowDisplayMenu()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Call the appropriate display function based on choice
	switch displayChoice {
	case "basic":
		interactive.DisplayBasicData(bars, symbol, timeframe)
	case "full":
		interactive.DisplayAdvancedData(bars, symbol, timeframe)
	case "analytics":
		interactive.DisplayAnalyticsData(bars, symbol, timeframe)
	case "all":
		interactive.DisplayBasicData(bars, symbol, timeframe)
		interactive.DisplayAdvancedData(bars, symbol, timeframe)
		interactive.DisplayAnalyticsData(bars, symbol, timeframe)
	}

	fmt.Printf("\n✅ Displayed %d bars for %s on %s timeframe\n", len(bars), symbol, timeframe)

	balanceChange := account.Equity.Sub(account.LastEquity)

	fmt.Println("Status:", resp.Status, balanceChange)

	// TODO: Add Obsidian export functionality here later
	fmt.Println("📝 Obsidian export functionality will be added in future updates")
}

func testIndicators(symbol string, days int, timeframe string) {
	log.Println("🧪 Testing RSI and ATR calculations...")

	// Fetch bars from Alpaca
	bars, err := datafeed.GetAlpacaBars(symbol, timeframe, days)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		return
	}
	log.Printf("✅ Fetched %d bars", len(bars))

	// Store them in database
	err = datafeed.StoreBarsWithAnalytics(symbol, timeframe, bars)
	if err != nil {
		log.Printf("Failed to store bars: %v", err)
		return
	}
	log.Println("✅ Stored bars in database")

	// Calculate limit based on timeframe (approximate bars per day)
	var calcLimit int
	switch timeframe {
	case "1Min":
		calcLimit = days * 390 // 390 minutes
	case "5Min":
		calcLimit = days * 78 // 390 / 5
	case "10Min":
		calcLimit = days * 39 // 390 / 10
	case "15Min":
		calcLimit = days * 26 // 390 / 15
	case "30Min":
		calcLimit = days * 13 // 390 / 30
	case "1Hour":
		calcLimit = days * 7 // ~6.5 hours per day
	case "1Day":
		calcLimit = days // 1 bar per day
	default:
		calcLimit = days * 10
	}

	// Calculate and store RSI
	err = strategy.CalculateAndStoreRSI(symbol, 14, timeframe, calcLimit)
	if err != nil {
		log.Printf("Failed to calculate RSI: %v", err)
		return
	}
	log.Println("✅ RSI calculation and storage successful!")

	// Calculate and store ATR
	err = strategy.CalculateAndStoreATR(symbol, 14, timeframe, calcLimit)
	if err != nil {
		log.Printf("Failed to calculate ATR: %v", err)
		return
	}
	log.Println("✅ ATR calculation and storage successful!")
}

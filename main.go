package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/strategy"
	"github.com/fazecat/mongelmaker/Internal/utils"
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
	utils.TestRetryLogic()

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
	// Test RSI calculation
	testRSI()

	// Get number of bars to fetch from user
	var numBars int
	fmt.Print("How many bars to display? ")
	fmt.Scan(&numBars)

	bars, err := interactive.FetchMarketData("AAPL", timeframe, numBars)
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
		interactive.DisplayBasicData(bars, "AAPL", timeframe)
	case "full":
		interactive.DisplayAdvancedData(bars, "AAPL", timeframe)
	case "analytics":
		interactive.DisplayAnalyticsData(bars, "AAPL", timeframe)
	case "all":
		interactive.DisplayBasicData(bars, "AAPL", timeframe)
		interactive.DisplayAdvancedData(bars, "AAPL", timeframe)
		interactive.DisplayAnalyticsData(bars, "AAPL", timeframe)
	}

	fmt.Printf("\n✅ Displayed %d bars for AAPL on %s timeframe\n", len(bars), timeframe)

	balanceChange := account.Equity.Sub(account.LastEquity)

	fmt.Println("Status:", resp.Status, balanceChange)

	// TODO: Add Obsidian export functionality here later
	fmt.Println("� Obsidian export functionality will be added in future updates")
}

func testRSI() {
	log.Println("🧪 Testing RSI calculation...")

	// 1. Fetch bars from Alpaca
	bars, err := datafeed.GetAlpacaBars("AAPL", "1Day", 50)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		return
	}
	log.Printf("✅ Fetched %d bars", len(bars))

	// 2. Store them in database
	err = datafeed.StoreBarsWithAnalytics("AAPL", "1Day", bars)
	if err != nil {
		log.Printf("Failed to store bars: %v", err)
		return
	}
	log.Println("✅ Stored bars in database")

	// 3. Calculate and store RSI
	err = strategy.CalculateAndStoreRSI("AAPL", 14)
	if err != nil {
		log.Printf("Failed to calculate RSI: %v", err)
		return
	}

	log.Println("✅ RSI calculation and storage successful!")
}

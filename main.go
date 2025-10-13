package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	err = datafeed.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer datafeed.CloseDatabase()

	// Test your new storage function
	testStorageFunction()

	apiKey := os.Getenv("ALPACA_API_KEY")
	secretKey := os.Getenv("ALPACA_SECRET_KEY")

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

	balanceChange := account.Equity.Sub(account.LastEquity)

	fmt.Println("Status:", resp.Status, balanceChange)
}

func testStorageFunction() {
	fmt.Println("\n🔄 Testing StoreBarsWithAnalytics function...")

	// Test with AAPL data
	symbol := "AAPL"
	fmt.Printf("📊 Fetching multi-timeframe data for %s...\n", symbol)

	// Fetch data using your existing function
	data, err := datafeed.FetchAllTimeframes(symbol, "1Day", 5)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		return
	}

	// Test your new StoreBarsWithAnalytics function
	fmt.Printf("💾 Storing %d bars of 1Min data with analytics...\n", len(data.OneMinData))
	err = datafeed.StoreBarsWithAnalytics(symbol, "1Min", data.OneMinData)
	if err != nil {
		log.Printf("Failed to store 1Min data: %v", err)
		return
	}

	fmt.Println("✅ Storage test completed!")
}

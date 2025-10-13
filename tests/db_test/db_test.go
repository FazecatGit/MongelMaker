package main

import (
	"fmt"
	"log"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load("../../.env") // Load from project root
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	err = datafeed.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer datafeed.CloseDatabase()

	fmt.Println("🔄 Starting database storage test...")

	// Test with AAPL data
	symbol := "AAPL"
	fmt.Printf("📊 Fetching multi-timeframe data for %s...\n", symbol)

	// Fetch data using your existing function
	data, err := datafeed.FetchAllTimeframes(symbol, "1Day", 5)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	// Test your new StoreBarsWithAnalytics function
	fmt.Printf("💾 Storing %d bars of 1Min data with analytics...\n", len(data.OneMinData))
	err = datafeed.StoreBarsWithAnalytics(symbol, "1Min", data.OneMinData)
	if err != nil {
		log.Fatalf("Failed to store 1Min data: %v", err)
	}

	fmt.Printf("💾 Storing %d bars of 5Min data with analytics...\n", len(data.FiveMinData))
	err = datafeed.StoreBarsWithAnalytics(symbol, "5Min", data.FiveMinData)
	if err != nil {
		log.Fatalf("Failed to store 5Min data: %v", err)
	}

	fmt.Printf("💾 Storing %d bars of 1Day data with analytics...\n", len(data.OneDayData))
	err = datafeed.StoreBarsWithAnalytics(symbol, "1Day", data.OneDayData)
	if err != nil {
		log.Fatalf("Failed to store 1Day data: %v", err)
	}

	// Query back some data to see the analytics
	fmt.Println("\n📈 Checking stored data with analytics...")
	testQuery := `
		SELECT symbol, timeframe, close_price, price_change, price_change_percent, timestamp 
		FROM historical_bars 
		WHERE symbol = $1 
		ORDER BY timestamp DESC 
		LIMIT 5
	`

	rows, err := datafeed.DB.Query(testQuery, symbol)
	if err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nRecent bars with analytics:")
	fmt.Println("Symbol | Timeframe | Close Price | Price Change | Change % | Timestamp")
	fmt.Println("-------|-----------|-------------|--------------|----------|----------")

	for rows.Next() {
		var symbol, timeframe, timestamp string
		var closePrice, priceChange, priceChangePercent float64

		err := rows.Scan(&symbol, &timeframe, &closePrice, &priceChange, &priceChangePercent, &timestamp)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		fmt.Printf("%-6s | %-9s | $%-10.2f | $%-11.2f | %-7.2f%% | %s\n",
			symbol, timeframe, closePrice, priceChange, priceChangePercent, timestamp[:19])
	}

	fmt.Println("\n🎉 Database storage test completed successfully!")
}

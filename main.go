package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/handlers"
	newsscraping "github.com/fazecat/mongelmaker/Internal/news_scraping"
	"github.com/fazecat/mongelmaker/Internal/utils"
	"github.com/fazecat/mongelmaker/Internal/utils/config"
	"github.com/fazecat/mongelmaker/Internal/utils/scanner"
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

	_, err = alpclient.GetAccount()
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Check market status
	cfg, _ := config.LoadConfig()
	status, isOpen := utils.CheckMarketStatus(time.Now(), cfg)
	fmt.Printf("📊 Market Status: %s (Open: %v)\n\n", status, isOpen)

	// Initialize Finnhub client for news fetching
	finnhubClient := newsscraping.NewFinnhubClient()
	_ = finnhubClient // TODO: Use if needed

	// Start background scanner goroutine
	ctx := context.Background()
	go startBackgroundScanner(ctx, cfg)

	// Main event loop
	for {
		fmt.Println("\n--- MongelMaker Menu ---")
		fmt.Println("1. Scan Watchlist")
		fmt.Println("2. Analyze Single Stock")
		fmt.Println("3. Run Screener")
		fmt.Println("4. View Watchlist")
		fmt.Println("5. Scout Symbols")
		fmt.Println("6. Exit")
		fmt.Print("Enter choice (1-6): ")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		switch choice {
		case 1:
			handlers.HandleScan(ctx, cfg, datafeed.Queries)
		case 2:
			handlers.HandleAnalyzeSingle(ctx, datafeed.Queries)
		case 3:
			handlers.HandleScreener(ctx, cfg, datafeed.Queries)
		case 4:
			handlers.HandleWatchlist(ctx, datafeed.Queries)
		case 5:
			handlers.HandleScout(ctx, cfg, datafeed.Queries)
		case 6:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

func startBackgroundScanner(ctx context.Context, cfg *config.Config) {
	log.Println("Background scanner started...")
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			log.Println("Background scanner stopped")
			return
		default:
			log.Println("Background scanner tick - checking for scans...")
			_, err := scanner.PerformScan(ctx, "default", cfg, datafeed.Queries)
			if err != nil {
				log.Printf("Background scan error: %v", err)
			} else {
				log.Println("Background scan completed successfully")
			}
		}
	}
}

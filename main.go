package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/export"
	newsscraping "github.com/fazecat/mongelmaker/Internal/news_scraping"
	"github.com/fazecat/mongelmaker/Internal/strategy"
	"github.com/fazecat/mongelmaker/Internal/utils"
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

	// Check market status
	cfg, _ := utils.LoadConfig()
	status, isOpen := utils.CheckMarketStatus(time.Now(), cfg)
	fmt.Printf("📊 Market Status: %s (Open: %v)\n\n", status, isOpen)

	// Initialize Finnhub client for news fetching
	finnhubClient := newsscraping.NewFinnhubClient()

	// Show main menu
	mainChoice, err := interactive.ShowMainMenu()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Show timeframe menu
	timeframe, err := interactive.ShowTimeframeMenu()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get number of bars to fetch/analyze from user
	var numBars int
	fmt.Print("\nHow many bars to fetch and analyze?: ")
	fmt.Scan(&numBars)

	var symbol string

	// Handle single vs screener path
	if mainChoice == "single" {
		fmt.Print("\nEnter stock symbol to analyze (e.g., AAPL, TSLA, MSFT): ")
		fmt.Scan(&symbol)
	} else if mainChoice == "screener" {
		fmt.Println("\n🕵️  Stock Screener")
		symbols := strategy.GetPopularStocks()[:50] // Top 50 stocks
		fmt.Printf("Screening %d popular stocks...\n", len(symbols))
		results, err := strategy.ScreenStocks(symbols, timeframe, numBars, strategy.DefaultScreenerCriteria(), nil)
		if err != nil {
			fmt.Printf("Screener failed: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("❌ No stocks matched the screener criteria. Try a different timeframe or number of bars.")
			return
		}

		// Display screener results
		fmt.Println("\n📈 Screener Results (Top Matches):")
		fmt.Println("Symbol | Score | RSI  | ATR  | Signals")
		fmt.Println("-------|-------|------|------|--------")
		for _, res := range results {
			rsiStr := "-"
			if res.RSI != nil {
				rsiStr = fmt.Sprintf("%.1f", *res.RSI)
			}
			atrStr := "-"
			if res.ATR != nil {
				atrStr = fmt.Sprintf("%.2f", *res.ATR)
			}
			signals := strings.Join(res.Signals, ", ")
			if len(signals) > 30 {
				signals = signals[:27] + "..."
			}
			fmt.Printf("%-6s | %5.1f | %4s | %4s | %s\n", res.Symbol, res.Score, rsiStr, atrStr, signals)
		}

		// Let user pick a stock from results
		symbol, err = interactive.PickStockFromResults(results)
		if err != nil {
			fmt.Println("Error picking stock:", err)
			return
		}
	}

	// Calculate indicators and fetch data for chosen symbol
	testIndicators(symbol, numBars, timeframe)

	// Fetch news for the chosen stock with Finnhub
	fmt.Printf("\n📰 Fetching news for %s...\n", symbol)
	news, err := finnhubClient.FetchNews(symbol, 5)
	if err != nil {
		fmt.Printf("⚠️  Could not fetch news: %v\n", err)
	} else if len(news) > 0 {
		fmt.Println("\n📰 Latest News Headlines:")
		for i, article := range news {
			sentimentEmoji := "➡️ "
			if article.Sentiment == newsscraping.Positive {
				sentimentEmoji = "📈"
			} else if article.Sentiment == newsscraping.Negative {
				sentimentEmoji = "📉"
			}

			impactStr := fmt.Sprintf("%.0f%%", article.Impact*100)

			fmt.Printf("\n%d. %s %s | Catalyst: %s | Impact: %s\n",
				i+1, sentimentEmoji, article.Sentiment, article.CatalystType, impactStr)
			fmt.Printf("   📰 %s\n", article.Headline)
			fmt.Printf("   🔗 %s\n", article.URL)
			fmt.Printf("   📅 %s\n", article.PublishedAt.Format("2006-01-02 15:04 MST"))
		}
	}

	bars, err := interactive.FetchMarketData(symbol, timeframe, numBars, "")
	if err != nil {
		fmt.Println("Error fetching market data:", err)
		return
	}

	// Show display menu
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
		interactive.DisplayAnalyticsData(bars, symbol, timeframe, timezone, datafeed.Queries)
	case "all":
		interactive.DisplayBasicData(bars, symbol, timeframe)
		interactive.DisplayAdvancedData(bars, symbol, timeframe)
		interactive.DisplayAnalyticsData(bars, symbol, timeframe, timezone, datafeed.Queries)
	case "vwap":
		interactive.DisplayVWAPAnalysis(bars, symbol, timeframe)
	case "export":
		records := interactive.PrepareExportData(bars, symbol, timezone)
		var format string
		fmt.Print("Enter export format (csv or json): ")
		fmt.Scan(&format)
		var filename string
		fmt.Print("Enter filename (e.g., data.csv): ")
		fmt.Scan(&filename)
		err := export.ExportData(format, filename, records)
		if err != nil {
			fmt.Printf("Export failed: %v\n", err)
		} else {
			fmt.Printf("✅ Exported to exported_data/%s\n", filename)
		}
	default:
		fmt.Println("Invalid display choice.")
	}

	fmt.Printf("\n✅ Analyzed %d bars for %s on %s timeframe\n", len(bars), symbol, timeframe)

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

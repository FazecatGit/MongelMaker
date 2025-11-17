package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	database "github.com/fazecat/mongelmaker/Internal/database/sqlc"
	"github.com/fazecat/mongelmaker/Internal/database/watchlist"
	newsscraping "github.com/fazecat/mongelmaker/Internal/news_scraping"
	"github.com/fazecat/mongelmaker/Internal/strategy"
	"github.com/fazecat/mongelmaker/Internal/utils/config"
	"github.com/fazecat/mongelmaker/Internal/utils/scanner"
	"github.com/fazecat/mongelmaker/Internal/utils/scoring"
	"github.com/fazecat/mongelmaker/interactive"
)

// clearInputBuffer clears any remaining input from stdin
func clearInputBuffer() {
	reader := bufio.NewReader(os.Stdin)
	for {
		r, _, err := reader.ReadRune()
		if err != nil || r == '\n' {
			break
		}
	}
}

func HandleScan(ctx context.Context, cfg *config.Config, q *database.Queries) {
	if len(cfg.Profiles) == 0 {
		fmt.Println("❌ No profiles configured")
		return
	}

	fmt.Println("\n📋 Available Profiles:")
	profiles := make([]string, 0)
	for name := range cfg.Profiles {
		profiles = append(profiles, name)
	}

	for i, profileName := range profiles {
		profile := cfg.Profiles[profileName]
		fmt.Printf("%d. %s (scan interval: %d days)\n", i+1, profileName, profile.ScanIntervalDays)
	}

	fmt.Print("Select profile (number): ")
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(profiles) {
		fmt.Println("❌ Invalid selection")
		return
	}

	selectedProfile := profiles[choice-1]

	fmt.Printf("🔄 Scanning profile: %s\n", selectedProfile)
	scannedCount, err := scanner.PerformScan(ctx, selectedProfile, cfg, q)
	if err != nil {
		fmt.Printf("❌ Scan failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Scan complete! Updated %d symbols\n", scannedCount)
}

func HandleAnalyzeSingle(ctx context.Context, q *database.Queries) {
	fmt.Print("Enter stock symbol (e.g., AAPL): ")
	var symbol string
	_, err := fmt.Scanln(&symbol)
	if err != nil || symbol == "" {
		fmt.Println("❌ Invalid symbol")
		return
	}

	timeframe, err := interactive.ShowTimeframeMenu()
	if err != nil {
		fmt.Println("❌ Invalid timeframe")
		return
	}

	fmt.Print("Enter number of bars (default 100): ")
	var numBars int
	_, err = fmt.Scanln(&numBars)
	if err != nil || numBars < 14 {
		numBars = 100
	}

	bars, err := interactive.FetchMarketData(symbol, timeframe, numBars, "")
	if err != nil {
		fmt.Printf("❌ Failed to fetch data: %v\n", err)
		return
	}

	err = datafeed.CalculateAndStoreRSI(symbol, bars)
	if err != nil {
		fmt.Printf("❌ Failed to calculate and store RSI: %v\n", err)
		return
	}

	err = datafeed.CalculateAndStoreATR(symbol, bars)
	if err != nil {
		fmt.Printf("❌ Failed to calculate and store ATR: %v\n", err)
		return
	}

	displayChoice, _ := interactive.ShowDisplayMenu()
	clearInputBuffer()

	switch displayChoice {
	case "basic":
		interactive.DisplayBasicData(bars, symbol, timeframe)
	case "full":
		interactive.DisplayAdvancedData(bars, symbol, timeframe)
	case "analytics":
		tz, _ := interactive.ShowTimezoneMenu()
		clearInputBuffer()
		interactive.DisplayAnalyticsData(bars, symbol, timeframe, tz, q)
		fmt.Println("\n--- Press Enter to continue ---")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	case "vwap":
		interactive.DisplayVWAPAnalysis(bars, symbol, timeframe)
	default:
		interactive.DisplayBasicData(bars, symbol, timeframe)
	}
}

func HandleScreener(ctx context.Context, cfg *config.Config, q *database.Queries) {
	symbols := strategy.GetPopularStocks()
	if len(symbols) == 0 {
		fmt.Println("❌ Could not get popular stocks")
		return
	}

	criteria := strategy.DefaultScreenerCriteria()

	fmt.Println("🔍 Screening stocks...")
	results, err := strategy.ScreenStocks(symbols, "1Day", 100, criteria, nil)
	if err != nil {
		fmt.Printf("❌ Screener failed: %v\n", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("📭 No stocks matched criteria")
		return
	}

	fmt.Printf("\n📊 Screening Results (%d total):\n", len(results))
	fmt.Println("==========================================")
	fmt.Println("# | Symbol | Score  | RSI    | ATR    | Signals                    | Analysis")
	fmt.Println("--|--------|--------|--------|--------|----------------------------|----------------------")

	for i, stock := range results {
		rsiStr := "  -   "
		if stock.RSI != nil {
			rsiStr = fmt.Sprintf("%6.2f", *stock.RSI)
		}

		atrStr := "  -   "
		if stock.ATR != nil {
			atrStr = fmt.Sprintf("%6.2f", *stock.ATR)
		}

		signalsStr := ""
		if len(stock.Signals) > 0 {
			for j, sig := range stock.Signals {
				if j > 0 {
					signalsStr += ", "
				}
				signalsStr += sig
			}
		} else {
			signalsStr = "-"
		}

		if len(signalsStr) > 26 {
			signalsStr = signalsStr[:23] + "..."
		}

		analysis := "---"
		if stock.RSI != nil {
			if *stock.RSI > 70 {
				analysis = "🔴 Overbought"
			} else if *stock.RSI < 30 {
				analysis = "🟢 Oversold"
			} else if *stock.RSI > 50 {
				analysis = "📈 Bullish"
			} else {
				analysis = "📉 Bearish"
			}
		}

		fmt.Printf("%2d| %s | %.2f | %s | %s | %-26s | %s\n",
			i+1, stock.Symbol, stock.Score, rsiStr, atrStr, signalsStr, analysis)
	}

	fmt.Print("\nSelect stock for details (or press Enter to skip): ")
	var choice int
	_, err = fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(results) {
		return
	}

	selectedStock := results[choice-1]

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("📊 Detailed Analysis: %s\n", selectedStock.Symbol)
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	fmt.Printf("🎯 Score: %.2f\n", selectedStock.Score)

	if selectedStock.RSI != nil {
		fmt.Printf("📈 RSI (14): %.2f", *selectedStock.RSI)
		if *selectedStock.RSI > 70 {
			fmt.Print(" 🔴 Overbought")
		} else if *selectedStock.RSI < 30 {
			fmt.Print(" 🟢 Oversold")
		}
		fmt.Println()
	}

	if selectedStock.ATR != nil {
		fmt.Printf("📊 ATR: %.2f", *selectedStock.ATR)
		if *selectedStock.ATR > 1.0 {
			fmt.Print(" ⚠️ High Volatility")
		}
		fmt.Println()
	}

	if len(selectedStock.Signals) > 0 {
		fmt.Println("\n🔔 Signals:")
		for _, sig := range selectedStock.Signals {
			fmt.Printf("   • %s\n", sig)
		}
	}

	if selectedStock.Recommendation != "" {
		fmt.Printf("\n📝 Recommendation: %s\n", selectedStock.Recommendation)
	}

	fmt.Println("\n📰 Fetching recent news...")
	finnhubClient := newsscraping.NewFinnhubClient()
	newsArticles, err := finnhubClient.FetchNews(selectedStock.Symbol, 5)
	if err != nil {
		fmt.Printf("⚠️ Could not fetch news: %v\n", err)
	} else if len(newsArticles) > 0 {
		fmt.Printf("\n📰 Recent News (%d articles):\n", len(newsArticles))
		fmt.Println(strings.Repeat("-", 80))
		for i, article := range newsArticles {
			sentimentIcon := "⚪"
			switch article.Sentiment {
			case newsscraping.Positive:
				sentimentIcon = "🟢"
			case newsscraping.Negative:
				sentimentIcon = "🔴"
			}

			catalystIcon := ""
			if article.CatalystType != newsscraping.NoCatalyst {
				catalystIcon = fmt.Sprintf(" [%s]", article.CatalystType)
			}

			fmt.Printf("\n%d. %s %s%s\n", i+1, sentimentIcon, article.Headline, catalystIcon)
			fmt.Printf("   🔗 %s\n", article.URL)
			fmt.Printf("   📅 %s\n", article.PublishedAt.Format("Jan 02, 2006 15:04"))
		}
		fmt.Println()
	} else {
		fmt.Println("📭 No recent news found")
	}

	fmt.Print("\n➕ Add to watchlist? (y/n): ")
	var addChoice string
	fmt.Scanln(&addChoice)

	if strings.ToLower(addChoice) == "y" {
		reason := "Added from screener"
		if selectedStock.Recommendation != "" {
			reason = fmt.Sprintf("Added from screener - %s", selectedStock.Recommendation)
			if len(reason) > 200 {
				reason = reason[:200]
			}
		}
		_, err = watchlist.AddToWatchlist(ctx, q, selectedStock.Symbol, "stock", selectedStock.Score, reason)
		if err != nil {
			fmt.Printf("❌ Failed to add to watchlist: %v\n", err)
			return
		}
		fmt.Printf("✅ Added %s to watchlist (Score: %.2f)\n", selectedStock.Symbol, selectedStock.Score)
	}

	fmt.Println("\n--- Press Enter to continue ---")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func HandleWatchlist(ctx context.Context, q *database.Queries) {
	watchlist, err := q.GetWatchlist(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to fetch watchlist: %v\n", err)
		return
	}
	fmt.Println("\n📋 Watchlist Menu:")
	fmt.Println("1. View Watchlist")
	fmt.Println("2. Exit")
	fmt.Print("Enter choice (number): ")

	var choice int
	_, err = fmt.Scanln(&choice)
	if err != nil {
		fmt.Println("❌ Invalid input")
		return
	}

	switch choice {
	case 1:
		if len(watchlist) == 0 {
			fmt.Println("📭 Watchlist is empty")
			return
		}
		fmt.Println("\n📊 Current Watchlist:")
		fmt.Println("Symbol | Score | Added Date | Last Updated | Category")
		fmt.Println("-------|-------|------------|--------------|---------")
		for _, item := range watchlist {
			addedStr := "N/A"
			if item.AddedDate.Valid {
				addedStr = item.AddedDate.Time.Format("2006-01-02")
			}
			updatedStr := "N/A"
			if item.LastUpdated.Valid {
				updatedStr = item.LastUpdated.Time.Format("2006-01-02")
			}
			fmt.Printf("%s | %.2f | %s | %s | %s\n", item.Symbol, item.Score, addedStr, updatedStr, scoring.ScoreCategory(float64(item.Score)))
		}
	case 2:
		return
	default:
		fmt.Println("❌ Invalid choice")
	}
}

func HandleScout(ctx context.Context, cfg *config.Config, q *database.Queries) {
	profiles := make([]string, 0)
	for name := range cfg.Profiles {
		profiles = append(profiles, name)
	}

	for i, profileName := range profiles {
		profile := cfg.Profiles[profileName]
		fmt.Printf("%d. %s (scan interval: %d days)\n", i+1, profileName, profile.ScanIntervalDays)
	}

	var minScore float64
	fmt.Print("Enter minimum score threshold (e.g., 0.0 to 100.0): ")
	_, err := fmt.Scanln(&minScore)
	if err != nil {
		fmt.Println("❌ Invalid input for minimum score threshold")
		return
	}

	fmt.Print("Select profile (number): ")
	var choice int
	_, err = fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(profiles) {
		fmt.Println("❌ Invalid selection")
		return
	}

	selectedProfile := profiles[choice-1]

	candidates, err := scanner.PerformProfileScan(ctx, selectedProfile, minScore)
	if err != nil {
		fmt.Printf("❌ Scout scan failed: %v\n", err)
		return
	}
	if len(candidates) == 0 {
		fmt.Println("📭 No candidates found above the score threshold")
		return
	}

	for i, candidate := range candidates {
		fmt.Printf("   %d. %s\n", i+1, candidate.Symbol)
		fmt.Printf("      Score: %.2f | Pattern: %s\n", candidate.Score, candidate.Analysis)

		fmt.Print("      ➕ Add to watchlist? (y/n): ")
		var addChoice string
		fmt.Scanln(&addChoice)
		if strings.ToLower(addChoice) == "y" {
			reason := fmt.Sprintf("Scouted - Pattern: %s", candidate.Analysis)
			_, err := watchlist.AddToWatchlist(ctx, q, candidate.Symbol, "stock", candidate.Score, reason)
			if err != nil {
				fmt.Printf("      ❌ Failed to add: %v\n", err)
				continue
			}
			fmt.Printf("      ✅ Added %s to watchlist\n", candidate.Symbol)
		}
	}

	fmt.Println("\n--- Press Enter to continue ---")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

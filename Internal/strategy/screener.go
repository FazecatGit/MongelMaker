package strategy

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	. "github.com/fazecat/mongelmaker/Internal/news_scraping"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

type ScreenerCriteria struct {
	MinOversoldRSI float64
	MaxRSI         float64
	MinATR         float64
	MinVolumeRatio float64
}

type StockScore struct {
	Symbol         string
	Score          float64
	Signals        []string
	RSI            *float64
	ATR            *float64
	NewsSentiment  SentimentScore
	NewsImpact     float64
	FinalSignal    CombinedSignal
	Recommendation string
}

func DefaultScreenerCriteria() ScreenerCriteria {
	return ScreenerCriteria{
		MinOversoldRSI: 35,  // RSI < 35 indicates oversold
		MaxRSI:         75,  // Avoid overbought
		MinATR:         0.1, // Very low volatility threshold
		MinVolumeRatio: 1.0, // Any volume above average
	}
}

// screens a list of symbols based on criteria
func ScreenStocks(symbols []string, timeframe string, numBars int, criteria ScreenerCriteria, newsStorage *NewsStorage) ([]StockScore, error) {
	var results []StockScore

	for _, symbol := range symbols {
		score, signals, rsi, atr, err := scoreStock(symbol, timeframe, numBars, criteria, newsStorage)
		if err != nil {
			log.Printf("Error screening %s: %v", symbol, err)
			continue
		}
		if score == 0 && len(signals) == 0 && rsi == nil && atr == nil {
			log.Printf("Skipping %s: no data available", symbol)
			continue
		}
		results = append(results, StockScore{
			Symbol:  symbol,
			Score:   score,
			Signals: signals,
			RSI:     rsi,
			ATR:     atr,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func scoreStock(symbol, timeframe string, numBars int, criteria ScreenerCriteria, newsStorage *NewsStorage) (score float64, signals []string, rsi, atr *float64, err error) {

	bars, err := datafeed.GetAlpacaBars(symbol, timeframe, numBars, "")
	if err != nil {
		return 0, nil, nil, nil, err
	}

	if len(bars) < 2 {
		return 0, nil, nil, nil, fmt.Errorf("insufficient data for %s (need 2 bars, got %d)", symbol, len(bars))
	}

	startTime := time.Now().AddDate(0, 0, -180)
	endTime := time.Now()

	if len(bars) > 0 {
		oldestTime, err := time.Parse(time.RFC3339, bars[len(bars)-1].Timestamp)
		if err == nil {
			startTime = oldestTime
		}
	}

	rsiMap, rsiErr := datafeed.FetchRSIByTimestampRange(symbol, startTime, endTime)
	if rsiErr != nil {
		log.Printf("RSI fetch failed for %s: %v (continuing with other signals)", symbol, rsiErr)
	} else if len(rsiMap) > 0 {
		rsi = findLatestValue(rsiMap)
	}

	atrMap, atrErr := datafeed.FetchATRByTimestampRange(symbol, startTime, endTime)
	if atrErr != nil {
		log.Printf("ATR fetch failed for %s: %v (continuing with other signals)", symbol, atrErr)
	} else if len(atrMap) > 0 {
		atr = findLatestValue(atrMap)
	}

	// Analyze latest candle
	latestBar := bars[0]
	volumes := make([]int64, len(bars))
	for i, bar := range bars {
		volumes[i] = bar.Volume
	}
	avgVol20 := utils.CalculateAvgVolume(volumes, 20)

	score = 0
	signals = []string{}

	if rsi != nil {
		if *rsi < criteria.MinOversoldRSI {
			score += 20
			signals = append(signals, fmt.Sprintf("RSI Oversold: %.2f", *rsi))
		} else if *rsi > criteria.MaxRSI {
			score -= 10
			signals = append(signals, fmt.Sprintf("RSI Overbought: %.2f", *rsi))
		} else {
			score += 5
		}
	}

	if atr != nil && *atr > criteria.MinATR {
		score += 10
		signals = append(signals, fmt.Sprintf("High Volatility ATR: %.2f", *atr))
	}

	if avgVol20 > 0 {
		volRatio := float64(latestBar.Volume) / avgVol20
		if volRatio > criteria.MinVolumeRatio {
			score += 15
			signals = append(signals, fmt.Sprintf("High Volume: %.1fx avg", volRatio))
		}
	}

	if newsStorage != nil {
		news, err := newsStorage.GetLatestNews(context.Background(), symbol, 1)
		if err == nil && len(news) > 0 && news[0].Sentiment == Positive {
			score += 10 // Boost score for positive sentiment
		}
	}

	whales := DetectWhales(symbol, bars)
	if len(whales) > 0 {
		for _, whale := range whales {
			// HIGH conviction whales get +5 score bonus
			if whale.Conviction == "HIGH" {
				score += 5
				signals = append(signals, fmt.Sprintf("🐋 Whale %s: Z=%.2f", whale.Direction, whale.ZScore))
			}
		}
	}

	support := FindSupport(bars)
	resistance := FindResistance(bars)

	currentPrice := latestBar.Close
	if currentPrice < support*1.01 {
		score += 15 // Strong buy signal
		signals = append(signals, fmt.Sprintf("Near Support: $%.2f", support))
	}
	if currentPrice > resistance*0.99 {
		score -= 10 // Sell signal
		signals = append(signals, fmt.Sprintf("Near Resistance: $%.2f", resistance))
	}

	combinedSignal := CalculateSignal(rsi, atr, bars, symbol, "")

	signals = append(signals, fmt.Sprintf("\n🎯 FINAL: %s", FormatSignal(combinedSignal)))

	return score, signals, rsi, atr, nil
}

func GetTradableAssets() ([]string, error) {
	client := datafeed.GetAlpacaClient()
	if client == nil {
		return nil, fmt.Errorf("alpaca client not initialized - call InitAlpacaClient() first")
	}

	assets, err := client.GetAssets(alpaca.GetAssetsRequest{
		Status: "active",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assets from Alpaca: %v", err)
	}

	symbols := make([]string, 0, len(assets))
	for _, asset := range assets {
		// Only include stocks, exclude options/crypto for now
		if asset.Class == "us_equity" && asset.Tradable {
			symbols = append(symbols, asset.Symbol)
		}
	}

	log.Printf("Fetched %d tradeable assets from Alpaca", len(symbols))
	return symbols, nil
}
func GetPopularStocks() []string {
	return []string{
		"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META", "NFLX", "BABA", "ORCL",
		"JPM", "BAC", "WFC", "C", "GS", "MS", "BLK", "AXP", "USB", "PNC",
		"JNJ", "PFE", "MRK", "ABT", "TMO", "DHR", "BMY", "LLY", "AMGN", "GILD",
		"XOM", "CVX", "COP", "EOG", "SLB", "HAL", "BKR", "OXY", "MPC", "PSX",
		"KO", "PEP", "MDLZ", "MO", "PM", "CL", "KMB", "GIS", "SYY", "HSY",
	}
}

// Helper to find latest value in map by timestamp keys
func findLatestValue(m map[string]float64) *float64 {
	if len(m) == 0 {
		return nil
	}
	var latestKey string
	var latestVal float64
	for k, v := range m {
		if latestKey == "" || k > latestKey {
			latestKey = k
			latestVal = v
		}
	}
	return &latestVal
}

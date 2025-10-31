package strategy

import (
	"fmt"
	"log"
	"sort"
	"time"

	datafeed "github.com/fazecat/mongelmaker/Internal/database"
	"github.com/fazecat/mongelmaker/Internal/utils"
)

type ScreenerCriteria struct {
	MinRSI            float64 // Minimum RSI for oversold (e.g., < 30)
	MaxRSI            float64 // Maximum RSI for overbought
	MinATR            float64 // Minimum ATR for volatility
	MinVolumeRatio    float64 // Minimum volume ratio vs average
	MinLowerWickRatio float64 // Minimum lower wick / body ratio for buying pressure
}

type StockScore struct {
	Symbol  string
	Score   float64
	Signals []string
	RSI     *float64
	ATR     *float64
}

func DefaultScreenerCriteria() ScreenerCriteria {
	return ScreenerCriteria{
		MinRSI:            0,   // No min, but we'll check < 30
		MaxRSI:            70,  // Avoid overbought
		MinATR:            0.5, // Some volatility
		MinVolumeRatio:    1.2, // 20% above average
		MinLowerWickRatio: 1.5, // Lower wick at least 1.5x body
	}
}

// ScreenStocks screens a list of symbols based on criteria
func ScreenStocks(symbols []string, timeframe string, numBars int, criteria ScreenerCriteria) ([]StockScore, error) {
	var results []StockScore

	for _, symbol := range symbols {
		score, signals, rsi, atr, err := scoreStock(symbol, timeframe, numBars, criteria)
		if err != nil {
			log.Printf("Error screening %s: %v", symbol, err)
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

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func scoreStock(symbol, timeframe string, numBars int, criteria ScreenerCriteria) (score float64, signals []string, rsi, atr *float64, err error) {
	// Fetch bars
	bars, err := datafeed.GetAlpacaBars(symbol, timeframe, numBars, "")
	if err != nil {
		return 0, nil, nil, nil, err
	}
	if len(bars) < 14 {
		return 0, nil, nil, nil, fmt.Errorf("insufficient data for %s", symbol)
	}

	// Fetch RSI
	startTime, err := time.Parse(time.RFC3339, bars[0].Timestamp)
	if err != nil {
		log.Printf("Failed to parse start time: %v", err)
	} else {
		endTime, err := time.Parse(time.RFC3339, bars[len(bars)-1].Timestamp)
		if err != nil {
			log.Printf("Failed to parse end time: %v", err)
		} else {
			rsiMap, err := datafeed.FetchRSIByTimestampRange(symbol, startTime, endTime)
			if err != nil {
				log.Printf("RSI fetch failed for %s: %v", symbol, err)
			} else if len(rsiMap) > 0 {
				rsi = findLatestValue(rsiMap)
			}
		}
	}

	// Fetch ATR
	startTime, err = time.Parse(time.RFC3339, bars[0].Timestamp)
	if err != nil {
		log.Printf("Failed to parse start time: %v", err)
	} else {
		endTime, err := time.Parse(time.RFC3339, bars[len(bars)-1].Timestamp)
		if err != nil {
			log.Printf("Failed to parse end time: %v", err)
		} else {
			atrMap, err := datafeed.FetchATRByTimestampRange(symbol, startTime, endTime)
			if err != nil {
				log.Printf("ATR fetch failed for %s: %v", symbol, err)
			} else if len(atrMap) > 0 {
				atr = findLatestValue(atrMap)
			}
		}
	}

	// Analyze latest candle
	latestBar := bars[len(bars)-1]
	avgVol20 := calculateAvgVolume(bars, 20)
	analysis, confidence := utils.PatternAnalyzeCandle(latestBar, atr, avgVol20, int64(latestBar.Volume))

	score = 0
	signals = []string{}

	// RSI scoring
	if rsi != nil {
		if *rsi < 30 {
			score += 20
			signals = append(signals, fmt.Sprintf("RSI Oversold: %.2f", *rsi))
		} else if *rsi > criteria.MaxRSI {
			score -= 10
			signals = append(signals, fmt.Sprintf("RSI Overbought: %.2f", *rsi))
		} else {
			score += 5
		}
	}

	// ATR scoring
	if atr != nil && *atr > criteria.MinATR {
		score += 10
		signals = append(signals, fmt.Sprintf("High Volatility ATR: %.2f", *atr))
	}

	// Volume scoring
	if avgVol20 > 0 {
		volRatio := float64(latestBar.Volume) / avgVol20
		if volRatio > criteria.MinVolumeRatio {
			score += 15
			signals = append(signals, fmt.Sprintf("High Volume: %.1fx avg", volRatio))
		}
	}

	// Pattern scoring
	if confidence > 0.7 {
		score += 10
		signals = append(signals, fmt.Sprintf("Strong Pattern: %s (%.0f%%)", analysis, confidence*100))
	} else if confidence > 0.5 {
		score += 5
		signals = append(signals, fmt.Sprintf("Pattern: %s (%.0f%%)", analysis, confidence*100))
	}

	return score, signals, rsi, atr, nil
}

func calculateAvgVolume(bars []datafeed.Bar, period int) float64 {
	if len(bars) < period {
		period = len(bars)
	}
	sum := 0.0
	for i := len(bars) - period; i < len(bars); i++ {
		sum += float64(bars[i].Volume)
	}
	return sum / float64(period)
}

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
func GetPopularStocks() []string {
	return []string{
		"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META", "NFLX", "BABA", "ORCL",
		"JPM", "BAC", "WFC", "C", "GS", "MS", "BLK", "AXP", "USB", "PNC",
		"JNJ", "PFE", "MRK", "ABT", "TMO", "DHR", "BMY", "LLY", "AMGN", "GILD",
		"XOM", "CVX", "COP", "EOG", "SLB", "HAL", "BKR", "OXY", "MPC", "PSX",
		"KO", "PEP", "MDLZ", "MO", "PM", "CL", "KMB", "GIS", "SYY", "HSY",
	}
}

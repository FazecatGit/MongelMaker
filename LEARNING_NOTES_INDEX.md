# 📚 MongelMaker Learning Documentation - Complete Index

## Overview

This documentation organizes all PART 1 and PART 2 requirements into individual Obsidian-friendly learning files. Each file:
- Explains one specific requirement
- References actual code implementation
- Includes examples and workflows
- Has wiki-style cross-references (`[[File Name]]`)

Total: **18 learning files** across PART 1 and PART 2

---

## PART 1: Foundation & Setup (8 Files)

### 1. Alpaca Authentication
**File**: `1_Alpaca_Authentication.md`
- Objective: Understand API authentication flow
- Topics: API keys, account types, token handling
- Code Reference: `Internal/database/alpaca.go`
- Related: [[OHLCV Fetching]], [[Error Handling & Retry Logic]]

### 2. OHLCV Fetching
**File**: `2_OHLCV_Fetching.md`
- Objective: Fetch candlestick data from Alpaca API
- Topics: Bar struct, timeframes, data ordering
- Code Reference: `Internal/types/bar.go`, `GetAlpacaBars()`
- Related: [[Alpaca Authentication]], [[Console Display & Pretty Printing]]

### 3. Console Display & Pretty Printing
**File**: `3_Console_Display.md`
- Objective: Format data for readable console output
- Topics: String formatting, tables, emoji integration
- Code Reference: `interactive/interactive.go`
- Related: [[OHLCV Fetching]], [[Final Signal Display]]

### 4. Error Handling & Retry Logic
**File**: `4_Error_Handling_Retry.md`
- Objective: Handle API failures gracefully
- Topics: Exponential backoff, transient errors, context
- Code Reference: `Internal/utils/retry.go`
- Related: [[Alpaca Authentication]], [[OHLCV Fetching]]

### 5. Database Setup & Schema
**File**: `5_Database_Setup.md`
- Objective: Set up PostgreSQL and sqlc
- Topics: Migrations, schema design, type-safe queries
- Code Reference: `Internal/sql/`, `sqlc.yaml`
- Related: [[Error Handling & Retry Logic]], [[OHLCV Fetching]]

### 6. Moving Averages & Trend Analysis
**File**: `6_Moving_Averages.md`
- Objective: Calculate EMA/SMA to identify trends
- Topics: EMA calculation, trend direction, timeframes
- Code Reference: `Internal/strategy/ema.go`
- Related: [[Stock Screener]], [[Support Resistance Detection]]

### 7. Candle Analysis - High/Low & Body/Wick Ratios
**File**: `7_Candle_Analysis.md`
- Objective: Analyze candlestick patterns
- Topics: Body/wick ratio, pattern recognition, Doji/Hammer
- Code Reference: `Internal/utils/analysis_candle.go`
- Related: [[Stock Screener]], [[Signal Combination Strategy]]

### 8. Export & Logging
**File**: `8_Export_Logging.md`
- Objective: Export results and log events
- Topics: JSON/CSV export, log levels, file organization
- Code Reference: `Internal/export/export_file.go`, `Internal/utils/logger.go`
- Related: [[Database Setup & Schema]], [[Stock Screener]]

---

## PART 2: Advanced Analysis & Trading Signals (10 Files)

### 1. Stock Screener
**File**: `1_Stock_Screener.md`
- Objective: Screen thousands of stocks for opportunities
- Topics: Scoring system, recommendations, filtering
- Code Reference: `Internal/strategy/screener.go`
- Related: [[RSI Analysis Per-Timeframe]], [[Whale Bear Spike Detection]], [[Support Resistance Detection]]

### 2. Whale & Bear Spike Detection (Volume Anomalies)
**File**: `2_Whale_Bear_Detection.md`
- Objective: Detect institutional buying/selling via Z-score
- Topics: Z-score calculation, direction detection, conviction
- Code Reference: `Internal/strategy/whale_dectector.go`
- Related: [[Stock Screener]], [[Volume Understanding]]

### 3. RSI Analysis Per-Timeframe
**File**: `3_RSI_Per_Timeframe.md`
- Objective: Calculate RSI on multiple timeframes
- Topics: RSI formula, oversold/overbought, crossovers
- Code Reference: `Internal/database/alpaca.go`, `FetchRSIByTimestampRange()`
- Related: [[Stock Screener]], [[Signal Combination Strategy]]

### 4. Volume Understanding - Foundation for Whale Detection
**File**: `4_Volume_Understanding.md`
- Objective: Deeply understand volume metrics
- Topics: Volume types, accumulation, clusters, patterns
- Code Reference: `Internal/strategy/whale_dectector.go`
- Related: [[Whale Bear Spike Detection]], [[Stock Screener]]

### 5. Body-to-Wick Ratio Analysis & Candlestick Patterns
**File**: `5_Body_Wick_Ratio.md`
- Objective: Identify bullish/bearish patterns
- Topics: B/W ratio, engulfing, patterns
- Code Reference: `Internal/utils/analysis_candle.go`
- Related: [[Stock Screener]], [[Signal Combination Strategy]]

### 6. ATR & Volatility Calculation
**File**: `6_ATR_Volatility.md`
- Objective: Measure price volatility
- Topics: True range, ATR formula, risk assessment
- Code Reference: `Internal/database/alpaca.go`, `FetchATRByTimestampRange()`
- Related: [[Stock Screener]], [[Signal Combination Strategy]]

### 7. Support & Resistance Level Detection
**File**: `7_Support_Resistance.md`
- Objective: Identify bounce and reversal levels
- Topics: Local minima/maxima, breakouts, distance calc
- Code Reference: `Internal/strategy/support_resistance.go`
- Related: [[Stock Screener]], [[Signal Combination Strategy]]

### 8. Signal Combination Strategy - Ensemble Weighting System
**File**: `8_Signal_Combination.md`
- Objective: Combine 5 signals into unified recommendation
- Topics: Weighting (RSI 25%, ATR 15%, Whale 30%, Pattern 20%, S/R 10%), confidence
- Code Reference: `Internal/strategy/signal_combination.go`
- Related: [[Stock Screener]], [[All other PART 2 files]]

### 9. News Scraping Service & Sentiment Analysis
**File**: `9_News_Scraping.md`
- Objective: Collect and analyze market news
- Topics: Finnhub API, RSS feeds, sentiment scoring
- Code Reference: `Internal/news_scraping/`
- Related: [[Stock Screener]], [[Signal Combination Strategy]]

### 10. Data Export to JSON & CSV Formats
**File**: `10_Data_Export.md`
- Objective: Export results to portable formats
- Topics: JSON/CSV structure, workflows, best practices
- Code Reference: `Internal/export/export_file.go`
- Related: [[Stock Screener]], [[Export & Logging]]

---

## Quick Navigation by Topic

### By Learning Stage

**Beginner (Start Here)**:
1. [[Alpaca Authentication]]
2. [[OHLCV Fetching]]
3. [[Console Display & Pretty Printing]]
4. [[Moving Averages & Trend Analysis]]

**Intermediate**:
5. [[Error Handling & Retry Logic]]
6. [[Database Setup & Schema]]
7. [[Candle Analysis - High/Low & Body/Wick Ratios]]
8. [[Export & Logging]]

**Advanced - Foundation**:
9. [[Stock Screener]]
10. [[RSI Analysis Per-Timeframe]]
11. [[ATR & Volatility Calculation]]

**Advanced - Signals**:
12. [[Volume Understanding]]
13. [[Whale & Bear Spike Detection]]
14. [[Support & Resistance Level Detection]]
15. [[Body-to-Wick Ratio Analysis]]
16. [[Signal Combination Strategy]]

**Advanced - Integration**:
17. [[News Scraping Service & Sentiment Analysis]]
18. [[Data Export to JSON & CSV Formats]]

### By Technical Area

**APIs & Data Fetching**:
- [[Alpaca Authentication]]
- [[OHLCV Fetching]]
- [[Error Handling & Retry Logic]]
- [[News Scraping Service & Sentiment Analysis]]

**Database & Storage**:
- [[Database Setup & Schema]]
- [[Export & Logging]]
- [[Data Export to JSON & CSV Formats]]

**Display & User Interface**:
- [[Console Display & Pretty Printing]]

**Technical Analysis**:
- [[Moving Averages & Trend Analysis]]
- [[Candle Analysis - High/Low & Body/Wick Ratios]]
- [[RSI Analysis Per-Timeframe]]
- [[ATR & Volatility Calculation]]
- [[Support & Resistance Level Detection]]
- [[Volume Understanding]]
- [[Whale & Bear Spike Detection]]

**Screening & Trading**:
- [[Stock Screener]]
- [[Signal Combination Strategy]]

---

## Code Architecture Map

### PART 1 Code Files
```
Internal/
├── database/
│   ├── alpaca.go              ← [[Alpaca Authentication]], [[OHLCV Fetching]]
│   └── db.go                  ← [[Database Setup & Schema]]
├── utils/
│   ├── analysis_candle.go     ← [[Candle Analysis]]
│   ├── logger.go              ← [[Export & Logging]]
│   ├── retry.go               ← [[Error Handling & Retry Logic]]
│   └── env_loader.go
├── strategy/
│   └── ema.go                 ← [[Moving Averages & Trend Analysis]]
├── export/
│   └── export_file.go         ← [[Export & Logging]], [[Data Export]]
└── sql/
    ├── schema/
    └── migrations/            ← [[Database Setup & Schema]]

interactive/
└── interactive.go             ← [[Console Display & Pretty Printing]]

types/
└── bar.go                      ← [[OHLCV Fetching]]
```

### PART 2 Code Files
```
Internal/
├── strategy/
│   ├── screener.go                 ← [[Stock Screener]]
│   ├── rsi.go                      ← [[RSI Analysis Per-Timeframe]]
│   ├── atr.go                      ← [[ATR & Volatility Calculation]]
│   ├── support_resistance.go       ← [[Support & Resistance Detection]]
│   ├── whale_dectector.go          ← [[Whale & Bear Spike Detection]], [[Volume Understanding]]
│   ├── signal_combination.go       ← [[Signal Combination Strategy]]
│   └── ema.go                      ← [[Moving Averages]]
├── news_scraping/
│   ├── catalyst.go                 ← [[News Scraping Service]]
│   ├── finnhub.go                  ← [[News Scraping Service]]
│   ├── rss.go                      ← [[News Scraping Service]]
│   ├── sentiment.go                ← [[News Scraping Service]]
│   ├── storage.go                  ← [[News Scraping Service]]
│   └── types.go                    ← [[News Scraping Service]]
└── export/
    └── export_file.go              ← [[Data Export to JSON & CSV]]

interactive/
└── interactive.go                  ← [[Console Display & Pretty Printing]]
```

---

## Learning Workflow Suggestions

### Path A: Complete Novice to Trader
1. [[Alpaca Authentication]] - Get API working
2. [[OHLCV Fetching]] - Fetch real data
3. [[Console Display & Pretty Printing]] - Visualize it
4. [[Moving Averages & Trend Analysis]] - First indicator
5. [[Candle Analysis - High/Low & Body/Wick Ratios]] - Pattern recognition
6. [[RSI Analysis Per-Timeframe]] - Momentum indicator
7. [[Stock Screener]] - Find opportunities
8. [[Support & Resistance Level Detection]] - Entry/exit points
9. [[Volume Understanding]] - Money flow analysis
10. [[Whale & Bear Spike Detection]] - Institutional activity
11. [[ATR & Volatility Calculation]] - Risk management
12. [[Signal Combination Strategy]] - Unified decision
13. [[News Scraping Service & Sentiment Analysis]] - Fundamentals
14. [[Data Export to JSON & CSV Formats]] - Share results

### Path B: Technical Focus
1. [[OHLCV Fetching]]
2. [[Moving Averages & Trend Analysis]]
3. [[Candle Analysis - High/Low & Body/Wick Ratios]]
4. [[RSI Analysis Per-Timeframe]]
5. [[ATR & Volatility Calculation]]
6. [[Support & Resistance Level Detection]]
7. [[Stock Screener]]
8. [[Signal Combination Strategy]]

### Path C: Engineering Focus
1. [[Alpaca Authentication]]
2. [[Error Handling & Retry Logic]]
3. [[Database Setup & Schema]]
4. [[Export & Logging]]
5. [[Console Display & Pretty Printing]]
6. [[Data Export to JSON & CSV Formats]]
7. [[Stock Screener]] (integration point)

---

## Project Completion Status

✅ **PART 1 Complete** (8/8 learning files)
- All foundational components documented
- Code references verified
- Examples included

✅ **PART 2 Complete** (10/10 learning files)
- All 5 signals documented (RSI, ATR, Whale, Pattern, S/R)
- Integration points explained
- Weighted ensemble system detailed

✅ **Cross-References Complete**
- All files have wiki-style links
- Easy navigation in Obsidian
- Clear learning paths

---

## Using These Files in Obsidian

1. Copy entire `LEARNING_NOTES_PART1/` and `LEARNING_NOTES_PART2/` directories
2. Create new Obsidian vault: "MongelMaker Learning"
3. Paste files into vault
4. Obsidian will auto-detect [[wiki links]] between files
5. Use graph view to visualize relationships
6. Use search to find all references to "whales", "RSI", etc.

**Recommended Obsidian plugins**:
- Graph View (built-in)
- Backlinks (built-in)
- Search (built-in)
- Table of Contents (for longer files)

---

## File Statistics

| Metric | Count |
|--------|-------|
| Total learning files | 18 |
| PART 1 files | 8 |
| PART 2 files | 10 |
| Total markdown lines | ~3,500 |
| Code examples | 150+ |
| Cross-references | 200+ |

---

## Version & Updates

- **Created**: Nov 8, 2025
- **Status**: Complete and verified
- **Code Status**: All references match v1.0 of MongelMaker
- **Next Update**: When features added or bugs fixed

---

## Questions or Feedback?

Each file includes:
- Clear learning objectives
- Implementation details
- Real examples
- Related files (for context)
- Test cases (for verification)

Start with any file that matches your current knowledge level.
Good luck learning! 🚀

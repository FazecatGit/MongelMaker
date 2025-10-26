-- +goose Up
-- Scout list for potential stocks
CREATE TABLE scout_list (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE NOT NULL,
    reason VARCHAR(255), -- "earnings_catalyst,oversold_rsi,sector_rotation"
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    trigger_price DECIMAL(10, 4),
    trigger_type VARCHAR(10), -- 'ABOVE' or 'BELOW'
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT
);

-- candle_daily_ for Bollinger Bands
CREATE TABLE candle_daily_bollinger (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    calculation_date DATE NOT NULL,
    time_period INTEGER NOT NULL, -- 30, 90, 365 days
    average_price DECIMAL(10, 4) NOT NULL,
    standard_deviation DECIMAL(10, 4) NOT NULL,
    upper_band DECIMAL(10, 4) NOT NULL,
    lower_band DECIMAL(10, 4) NOT NULL,
    current_price DECIMAL(10, 4) NOT NULL,
    deviation_from_avg DECIMAL(10, 4) NOT NULL,
    UNIQUE(symbol, calculation_date, time_period)
);

-- RSI calculations table
CREATE TABLE rsi_calculation(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    calculation_date DATE NOT NULL,
    rsi_value DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, calculation_date)
);

-- Daily candles table
CREATE TABLE candles_daily(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    candle_date DATE NOT NULL,
    open_price DECIMAL(10, 4) NOT NULL,
    high_price DECIMAL(10, 4) NOT NULL,
    low_price DECIMAL(10, 4) NOT NULL,
    close_price DECIMAL(10, 4) NOT NULL,
    volume BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, candle_date)
);


-- 12-hour candles table
CREATE TABLE candles_12h(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    candle_timestamp TIMESTAMP NOT NULL,
    open_price DECIMAL(10, 4) NOT NULL,
    high_price DECIMAL(10, 4) NOT NULL,
    low_price DECIMAL(10, 4) NOT NULL,
    close_price DECIMAL(10, 4) NOT NULL,
    volume BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, candle_timestamp)
);

-- Market news articles table
CREATE TABLE news_articles(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    headline TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    published_at TIMESTAMP NOT NULL,
    source VARCHAR(100),
    sentiment VARCHAR(10), -- 'POSITIVE', 'NEGATIVE', 'NEUTRAL'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Whale alerts table
CREATE TABLE whale_alerts(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    alert_type VARCHAR(50) NOT NULL, -- 'LARGE_BUY', 'LARGE_SELL', etc.
    amount DECIMAL(20, 4) NOT NULL,
    price DECIMAL(10, 4) NOT NULL,
    alert_timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Support and resistance levels table
CREATE TABLE support_levels(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    level_price DECIMAL(10, 4) NOT NULL,
    level_type VARCHAR(10) NOT NULL, -- 'SUPPORT' or 'RESISTANCE'
    detected_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance indexes for fast queries
CREATE INDEX idx_rsi_symbol_date ON rsi_calculation(symbol, calculation_date);
CREATE INDEX idx_rsi_date_value ON rsi_calculation(calculation_date, rsi_value);

CREATE INDEX idx_bollinger_symbol_period ON candle_daily_bollinger(symbol, time_period, calculation_date);

CREATE INDEX idx_candles_daily_symbol ON candles_daily(symbol, candle_date);
CREATE INDEX idx_candles_12h_symbol ON candles_12h(symbol, candle_timestamp);

CREATE INDEX idx_news_symbol_date ON news_articles(symbol, published_at);
CREATE INDEX idx_news_date ON news_articles(published_at DESC);

CREATE INDEX idx_whale_symbol_time ON whale_alerts(symbol, alert_timestamp);
CREATE INDEX idx_whale_time ON whale_alerts(alert_timestamp DESC);

CREATE INDEX idx_support_symbol_type ON support_levels(symbol, level_type);

-- +goose Down 
DROP TABLE IF EXISTS scout_list;
DROP TABLE IF EXISTS candle_daily_bollinger;
DROP TABLE IF EXISTS rsi_calculation;
DROP TABLE IF EXISTS candles_12h;
DROP TABLE IF EXISTS candles_daily;
DROP TABLE IF EXISTS news_articles;
DROP TABLE IF EXISTS whale_alerts;
DROP TABLE IF EXISTS support_levels;
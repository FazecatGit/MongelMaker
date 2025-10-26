-- MongelMaker Trading Bot Database Schema - UP ONLY

-- Historical OHLCV data table
CREATE TABLE historical_bars (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    timeframe VARCHAR(10) NOT NULL, -- '1Min', '5Min', '1Day', etc.
    timestamp TIMESTAMP NOT NULL,
    open_price DECIMAL(10, 4) NOT NULL,
    high_price DECIMAL(10, 4) NOT NULL,
    low_price DECIMAL(10, 4) NOT NULL,
    close_price DECIMAL(10, 4) NOT NULL,
    volume BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure no duplicate bars for same symbol/timeframe/timestamp
    UNIQUE(symbol, timeframe, timestamp)
);

-- Trading signals table
CREATE TABLE signals (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    signal_type VARCHAR(10) NOT NULL, -- 'BUY', 'SELL', 'HOLD'
    current_price DECIMAL(10, 4) NOT NULL,
    sma_value DECIMAL(10, 4),
    confidence DECIMAL(3, 2), -- 0.00 to 1.00
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    executed BOOLEAN DEFAULT FALSE
);

-- Executed trades table
CREATE TABLE trades (
    id SERIAL PRIMARY KEY,
    signal_id INTEGER REFERENCES signals(id),
    symbol VARCHAR(10) NOT NULL,
    side VARCHAR(4) NOT NULL, -- 'BUY' or 'SELL'
    quantity DECIMAL(10, 4) NOT NULL,
    price DECIMAL(10, 4) NOT NULL,
    total_value DECIMAL(12, 4) NOT NULL,
    commission DECIMAL(8, 4) DEFAULT 0,
    alpaca_order_id VARCHAR(50), -- Alpaca's order ID
    status VARCHAR(20) DEFAULT 'PENDING', -- 'PENDING', 'FILLED', 'CANCELLED'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    filled_at TIMESTAMP
);

-- Current positions table
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE NOT NULL,
    quantity DECIMAL(10, 4) NOT NULL,
    avg_entry_price DECIMAL(10, 4) NOT NULL,
    current_price DECIMAL(10, 4),
    market_value DECIMAL(12, 4),
    unrealized_pnl DECIMAL(12, 4),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Portfolio performance history
CREATE TABLE portfolio_history (
    id SERIAL PRIMARY KEY,
    total_equity DECIMAL(12, 4) NOT NULL,
    cash_balance DECIMAL(12, 4) NOT NULL,
    positions_value DECIMAL(12, 4) NOT NULL,
    day_change DECIMAL(12, 4),
    total_return DECIMAL(12, 4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scout list for potential stocks
CREATE TABLE scout_list (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE NOT NULL,
    reason VARCHAR(255), -- "earnings_catalyst,oversold_rsi,sector_rotation"
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    trigger_price DECIMAL(10, 4),
    trigger_type VARCHAR(10), -- 'ABOVE' or 'BELOW'
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
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
    deviation_from_avg DECIMAL(10, 4) NOT NULL
);

CREATE TABLE rsi_calculation(
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    calculation_date DATE NOT NULL,
    rsi_value DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_historical_bars_symbol_timeframe ON historical_bars(symbol, timeframe);
CREATE INDEX idx_historical_bars_timestamp ON historical_bars(timestamp);
CREATE INDEX idx_signals_symbol_created ON signals(symbol, created_at);
CREATE INDEX idx_trades_symbol_created ON trades(symbol, created_at);
CREATE INDEX idx_positions_symbol ON positions(symbol);
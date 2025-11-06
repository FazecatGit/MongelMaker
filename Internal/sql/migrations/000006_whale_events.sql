-- +goose Up
CREATE TABLE IF NOT EXISTS whale_events (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    direction TEXT NOT NULL,
    volume BIGINT NOT NULL,
    z_score DECIMAL NOT NULL,
    close_price DECIMAL NOT NULL,
    price_change DECIMAL,
    conviction TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_whale_symbol_time ON whale_events(symbol, timestamp DESC);

-- +goose Down
DROP TABLE IF EXISTS whale_events;
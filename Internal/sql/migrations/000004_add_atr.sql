-- +goose Up
CREATE TABLE atr_calculation (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    calculation_date DATE NOT NULL,
    atr_value DECIMAL(10, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, calculation_date)
);

CREATE INDEX idx_atr_symbol_date ON atr_calculation(symbol, calculation_date);

-- +goose down
DROP INDEX idx_atr_symbol_date;
DROP TABLE atr_calculation;
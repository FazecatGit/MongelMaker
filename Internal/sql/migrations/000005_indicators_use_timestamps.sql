-- +goose Up
-- Modify rsi_calculation to use timestamp instead of date
-- SQLite: recreate table with new schema (can't alter columns)
ALTER TABLE rsi_calculation RENAME TO rsi_calculation_old;

CREATE TABLE rsi_calculation (
  symbol TEXT NOT NULL,
  calculation_timestamp TIMESTAMP NOT NULL,
  rsi_value REAL NOT NULL,
  UNIQUE (symbol, calculation_timestamp)
);

INSERT INTO rsi_calculation SELECT symbol, calculation_date, rsi_value FROM rsi_calculation_old;
DROP TABLE rsi_calculation_old;

-- Same for atr_calculation
ALTER TABLE atr_calculation RENAME TO atr_calculation_old;

CREATE TABLE atr_calculation (
  symbol TEXT NOT NULL,
  calculation_timestamp TIMESTAMP NOT NULL,
  atr_value REAL NOT NULL,
  UNIQUE (symbol, calculation_timestamp)
);

INSERT INTO atr_calculation SELECT symbol, calculation_date, atr_value FROM atr_calculation_old;
DROP TABLE atr_calculation_old;

-- +goose Down
-- Revert changes
ALTER TABLE rsi_calculation 
  DROP CONSTRAINT rsi_calculation_symbol_timestamp_key;

ALTER TABLE rsi_calculation 
  RENAME COLUMN calculation_timestamp TO calculation_date;

ALTER TABLE rsi_calculation 
  ALTER COLUMN calculation_date TYPE DATE;

ALTER TABLE rsi_calculation 
  ADD CONSTRAINT rsi_calculation_symbol_calculation_date_key 
  UNIQUE (symbol, calculation_date);

ALTER TABLE atr_calculation 
  DROP CONSTRAINT atr_calculation_symbol_timestamp_key;

ALTER TABLE atr_calculation 
  RENAME COLUMN calculation_timestamp TO calculation_date;

ALTER TABLE atr_calculation 
  ALTER COLUMN calculation_date TYPE DATE;

ALTER TABLE atr_calculation 
  ADD CONSTRAINT atr_calculation_symbol_calculation_date_key 
  UNIQUE (symbol, calculation_date);

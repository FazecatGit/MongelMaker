-- +goose Up
-- Modify rsi_calculation to use timestamp instead of date
ALTER TABLE rsi_calculation 
  DROP CONSTRAINT rsi_calculation_symbol_calculation_date_key;

ALTER TABLE rsi_calculation 
  RENAME COLUMN calculation_date TO calculation_timestamp;

ALTER TABLE rsi_calculation 
  ALTER COLUMN calculation_timestamp TYPE TIMESTAMP;

ALTER TABLE rsi_calculation 
  ADD CONSTRAINT rsi_calculation_symbol_timestamp_key 
  UNIQUE (symbol, calculation_timestamp);

-- Modify atr_calculation to use timestamp instead of date
ALTER TABLE atr_calculation 
  DROP CONSTRAINT atr_calculation_symbol_calculation_date_key;

ALTER TABLE atr_calculation 
  RENAME COLUMN calculation_date TO calculation_timestamp;

ALTER TABLE atr_calculation 
  ALTER COLUMN calculation_timestamp TYPE TIMESTAMP;

ALTER TABLE atr_calculation 
  ADD CONSTRAINT atr_calculation_symbol_timestamp_key 
  UNIQUE (symbol, calculation_timestamp);

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

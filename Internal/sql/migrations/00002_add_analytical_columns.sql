-- +goose Up
-- Add analytical columns to historical_bars table
ALTER TABLE historical_bars ADD COLUMN price_change DECIMAL(10, 4);
ALTER TABLE historical_bars ADD COLUMN price_change_percent DECIMAL(8, 4);

-- +goose Down  
-- Remove analytical columns from historical_bars table
ALTER TABLE historical_bars DROP COLUMN IF EXISTS price_change;
ALTER TABLE historical_bars DROP COLUMN IF EXISTS price_change_percent;

-- name: GetClosingPrices :many
SELECT close_price, timestamp
FROM historical_bars
WHERE symbol = $1 
  AND timeframe = '1Day'
ORDER BY timestamp DESC
LIMIT $2;

-- name: SaveRSI :exec
INSERT INTO rsi_calculation (symbol, calculation_date, rsi_value)
VALUES ($1, $2, $3)
ON CONFLICT (symbol, calculation_date)
DO UPDATE SET rsi_value = EXCLUDED.rsi_value;

-- name: GetLatestRSI :one
SELECT rsi_value, calculation_date
FROM rsi_calculation
WHERE symbol = $1
ORDER BY calculation_date DESC
LIMIT 1;
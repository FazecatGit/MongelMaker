-- name: GetClosingPrices :many
SELECT close_price, timestamp
FROM historical_bars
WHERE symbol = $1 
  AND timeframe = $2
ORDER BY timestamp ASC
LIMIT $3;

-- name: SaveRSI :exec
INSERT INTO rsi_calculation (symbol, calculation_timestamp, rsi_value)
VALUES ($1, $2, $3)
ON CONFLICT (symbol, calculation_timestamp)
DO UPDATE SET rsi_value = EXCLUDED.rsi_value;

-- name: GetLatestRSI :one
SELECT rsi_value, calculation_timestamp
FROM rsi_calculation
WHERE symbol = $1
ORDER BY calculation_timestamp DESC
LIMIT 1;

-- name: GetATRPrices :many
SELECT high_price, low_price, close_price, timestamp
FROM historical_bars
WHERE symbol = $1 
  AND timeframe = $2
ORDER BY timestamp ASC
LIMIT $3;

-- name: GetATR :one
SELECT atr_value, calculation_timestamp
FROM atr_calculation
WHERE symbol = $1
ORDER BY calculation_timestamp DESC
LIMIT 1;

-- name: SaveATR :exec
INSERT INTO atr_calculation (symbol, calculation_timestamp, atr_value)
VALUES ($1, $2, $3)
ON CONFLICT (symbol, calculation_timestamp)
DO UPDATE SET atr_value = EXCLUDED.atr_value;

-- name: GetRSIForDateRange :many
SELECT calculation_timestamp, rsi_value
FROM rsi_calculation
WHERE symbol = $1
ORDER BY calculation_timestamp DESC
LIMIT $2;

-- name: GetATRForDateRange :many
SELECT calculation_timestamp, atr_value
FROM atr_calculation
WHERE symbol = $1
ORDER BY calculation_timestamp DESC
LIMIT $2;

-- name: GetRSIByTimestampRange :many
SELECT calculation_timestamp, rsi_value
FROM rsi_calculation
WHERE symbol = $1
  AND calculation_timestamp >= $2
  AND calculation_timestamp <= $3
ORDER BY calculation_timestamp ASC;

-- name: GetATRByTimestampRange :many
SELECT calculation_timestamp, atr_value
FROM atr_calculation
WHERE symbol = $1
  AND calculation_timestamp >= $2
  AND calculation_timestamp <= $3
ORDER BY calculation_timestamp ASC;
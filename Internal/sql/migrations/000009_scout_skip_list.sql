-- +goose Up
Create TABLE scout_skip_list (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    profile_name TEXT NOT NULL,
    asset_type TEXT NOT NULL,
    reason TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    recheck_after TIMESTAMP NOT NULL,  -- needs to be 2 days from symbol made for now
    CONSTRAINT unique_symbol_profile UNIQUE (symbol, profile_name)
);
-- +goose Down
DROP TABLE scout_skip_list;
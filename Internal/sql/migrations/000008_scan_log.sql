-- +goose Up
CREATE TABLE scan_log (
    id SERIAL PRIMARY KEY,
    profile_name VARCHAR(50) NOT NULL UNIQUE,
    last_scan_timestamp TIMESTAMP NOT NULL,
    next_scan_due TIMESTAMP NOT NULL,
    symbols_scanned INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE scan_log;

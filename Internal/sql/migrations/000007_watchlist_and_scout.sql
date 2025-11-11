-- +goose Up
-- Create watchlist table (main tracked candidates)
CREATE TABLE watchlist (
  id INTEGER PRIMARY KEY,
  symbol TEXT NOT NULL UNIQUE,
  asset_type TEXT NOT NULL,
  score REAL NOT NULL,
  reason TEXT,  
  added_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  status TEXT DEFAULT 'active'  -- 'active' or 'archived'
);

-- Create history table (track score changes over time)
CREATE TABLE watchlist_history (
  id INTEGER PRIMARY KEY,
  watchlist_id INTEGER NOT NULL,
  old_score REAL,
  new_score REAL NOT NULL,
  analysis_data TEXT, -- JSON
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(watchlist_id) REFERENCES watchlist(id)
);

-- Create skip table (rejected candidates with recheck date)
CREATE TABLE skip_backlog (
  id INTEGER PRIMARY KEY,
  symbol TEXT NOT NULL UNIQUE,
  asset_type TEXT NOT NULL,
  reason TEXT,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  recheck_after TIMESTAMP NOT NULL  -- when can we reconsider (default: 7 days from now)
);

-- +goose Down
DROP TABLE IF EXISTS skip_backlog;
DROP TABLE IF EXISTS watchlist_history;
DROP TABLE IF EXISTS watchlist;
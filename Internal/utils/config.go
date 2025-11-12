package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Global struct {
		MarketHours struct {
			RegularOpen    string `yaml:"regular_open"`
			RegularClose   string `yaml:"regular_close"`
			PremarketOpen  string `yaml:"premarket_open"`
			AfterhourClose string `yaml:"afterhours_close"`
			Timezone       string `yaml:"timezone"`
		} `yaml:"market_hours"`
		LiquidityMinimumUSD int `yaml:"liquidity_minimum_usd"`
	} `yaml:"global"`

	Notifications struct {
		Channels struct {
			Console bool `yaml:"console"`
			FileLog bool `yaml:"file_log"`
			Discord bool `yaml:"discord"`
		} `yaml:"channels"`
		BatchDigestTime string `yaml:"batch_digest_time"`
	} `yaml:"notifications"`

	Archive struct {
		DaysBeforeArchive    int `yaml:"days_before_archive"`
		RecheckSkipAfterDays int `yaml:"recheck_skip_after_days"`
	} `yaml:"archive"`

	Profiles map[string]struct {
		Threshold        float64 `yaml:"threshold"`
		ScanIntervalDays int     `yaml:"scan_interval_days"`
	} `yaml:"profiles"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

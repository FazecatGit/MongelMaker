package newsscraping

import "time"

type SentimentScore string

const (
	Positive SentimentScore = "POSITIVE"
	Negative SentimentScore = "NEGATIVE"
	Neutral  SentimentScore = "NEUTRAL"
)

type CatalystType string

const (
	Earnings    CatalystType = "EARNINGS"
	Acquisition CatalystType = "ACQUISITION"
	Regulatory  CatalystType = "REGULATORY"
	Leadership  CatalystType = "LEADERSHIP"
	Market      CatalystType = "MARKET"
	Technical   CatalystType = "TECHNICAL"
	NoCatalyst  CatalystType = "NO_CATALYST"
)

type NewsArticle struct {
	ID           int64
	Symbol       string
	Headline     string
	URL          string
	PublishedAt  time.Time
	Source       string
	Sentiment    SentimentScore
	CatalystType CatalystType
	Impact       float64
	CreatedAt    time.Time
}

type NewsScraper interface {
	FetchNews(symbol string, limit int) ([]NewsArticle, error)
	Name() string
}

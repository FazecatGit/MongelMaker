package newsscraping

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSSClient struct {
	feeds      map[string]string
	httpClient *http.Client
}

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Items []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func NewRSSClinet() *RSSClient {
	return &RSSClient{
		feeds: map[string]string{
			"AAPL":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=AAPL",
			"MSFT":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=MSFT",
			"GOOGL": "https://feeds.finance.yahoo.com/rss/2.0/headline?s=GOOGL",
			"AMZN":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=AMZN",
			"TSLA":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=TSLA",
			"META":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=META",
			"NVDA":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=NVDA",
			"NFLX":  "https://feeds.finance.yahoo.com/rss/2.0/headline?s=NFLX",
			"JPM":   "https://feeds.finance.yahoo.com/rss/2.0/headline?s=JPM",
			"BAC":   "https://feeds.finance.yahoo.com/rss/2.0/headline?s=BAC",
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *RSSClient) FetchNews(symbol string, limit int) ([]NewsArticle, error) {
	url, exists := c.feeds[symbol]
	if !exists {
		return nil, fmt.Errorf("no RSS feed for symbol: %s", symbol)
	}
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %v", err)
	}

	var articles []NewsArticle
	for i, item := range feed.Channel.Items {
		if i >= limit {
			break
		}

		pubTime, _ := time.Parse(time.RFC1123Z, item.PubDate)
		if time.Since(pubTime) > 7*24*time.Hour {
			continue // Skip articles older than 7 days
		}

		article := NewsArticle{
			Symbol:       symbol,
			Headline:     item.Title,
			URL:          item.Link,
			PublishedAt:  pubTime,
			Source:       "Yahoo Finance RSS",
			Sentiment:    Neutral,
			CatalystType: NoCatalyst,
			Impact:       0.0,
			CreatedAt:    time.Now(),
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (c *RSSClient) Name() string {
	return "Yahoo Finance RSS"
}

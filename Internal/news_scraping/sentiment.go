package newsscraping

import "strings"

type SentimentAnalyzer struct {
	positiveWords map[string]float64
	negativeWords map[string]float64
}

func NewSentimentAnalyzer() *SentimentAnalyzer {
	return &SentimentAnalyzer{
		positiveWords: map[string]float64{
			"surge": 1.0, "soar": 1.0, "rally": 0.9, "beat": 0.8,
			"gain": 0.7, "profit": 0.8, "growth": 0.8, "bullish": 0.9,
			"upgrade": 0.9, "strong": 0.7, "jump": 0.8, "recover": 0.7,
			"outperform": 0.8, "breakout": 0.8, "momentum": 0.7,
		},
		negativeWords: map[string]float64{
			"crash": 1.0, "plunge": 1.0, "miss": 0.8, "fall": 0.7,
			"loss": 0.8, "bankruptcy": 1.0, "bearish": 0.9, "downgrade": 0.9,
			"weak": 0.7, "decline": 0.6, "warning": 0.7, "risk": 0.5,
			"underperform": 0.8, "slump": 0.8,
		},
	}
}

func (sa *SentimentAnalyzer) Analyze(text string) (SentimentScore, float64) {
	text = strings.ToLower(text)
	words := strings.Fields(text)

	var score float64
	var matches int

	for _, word := range words {
		word = strings.Trim(word, ".,!?\"'()[]{}:;")

		if val, exists := sa.positiveWords[word]; exists {
			score += val
			matches++
		} else if val, exists := sa.negativeWords[word]; exists {
			score -= val
			matches++
		}
	}

	if matches > 0 {
		score /= float64(matches)
	}
	sentiment := Neutral
	if score > 0.1 {
		sentiment = Positive
	} else if score < -0.1 {
		sentiment = Negative
	}
	return sentiment, score

}

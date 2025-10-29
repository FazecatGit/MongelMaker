package utils

import (
	"math"
)

type Candlestick struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
}
type CandleMetrics struct {
	Body        float64
	Range       float64
	BodyPct     float64
	UpperWick   float64
	LowerWick   float64
	UpperPct    float64
	LowerPct    float64
	BodyToUpper float64
	BodyToLower float64
}

func CalculateBodyWickRatios(candle Candlestick) (bodyToUpperRatio, bodyToLowerRatio float64) {
	body := math.Abs(candle.Close - candle.Open)
	upperWick := candle.High - math.Max(candle.Open, candle.Close)
	lowerWick := math.Min(candle.Open, candle.Close) - candle.Low

	bodyToUpperRatio = 0
	bodyToLowerRatio = 0

	if upperWick != 0 {
		bodyToUpperRatio = body / upperWick
	}
	if lowerWick != 0 {
		bodyToLowerRatio = body / lowerWick
	}

	return
}

func AnalyzeCandlestick(candle Candlestick) (map[string]float64, map[string]string) {
	body := math.Abs(candle.Close - candle.Open)
	rangeVal := candle.High - candle.Low
	bodyPct := 0.0
	if rangeVal != 0 {
		bodyPct = (body / rangeVal) * 100
	}
	upperWick := candle.High - math.Max(candle.Open, candle.Close)
	lowerWick := math.Min(candle.Open, candle.Close) - candle.Low
	upperPct := 0.0
	lowerPct := 0.0
	if rangeVal != 0 {
		upperPct = (upperWick / rangeVal) * 100
		lowerPct = (lowerWick / rangeVal) * 100
	}
	bodyToUpper, bodyToLower := CalculateBodyWickRatios(candle)
	metrics := map[string]float64{
		"Body":        body,
		"Range":       rangeVal,
		"BodyPct":     bodyPct,
		"UpperWick":   upperWick,
		"LowerWick":   lowerWick,
		"UpperPct":    upperPct,
		"LowerPct":    lowerPct,
		"BodyToUpper": bodyToUpper,
		"BodyToLower": bodyToLower,
	}
	analysis := "Neutral"
	if bodyPct < 10 {
		if upperPct > 70 {
			analysis = "Bearish Rejection"
		} else if lowerPct > 70 {
			analysis = "Bullish Rejection"
		} else {
			analysis = "Doji (indecision)"
		}
	} else if candle.Close > candle.Open {
		if bodyPct > 60 {
			analysis = "Strong Bullish"
		} else {
			analysis = "Bullish"
		}
	} else {
		if bodyPct > 60 {
			analysis = "Strong Bearish"
		} else {
			analysis = "Bearish"
		}
	}

	stringResults := map[string]string{
		"Analysis": analysis,
	}

	return metrics, stringResults
}

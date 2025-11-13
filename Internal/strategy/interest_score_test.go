package strategy

import (
	"testing"
)

func TestCalculateInterestScore_NeutralConditions(t *testing.T) {
	input := ScoringInput{
		CurrentPrice: 100,
		VWAPPrice:    100,
		ATRValue:     1.0,
		RSIValue:     50,
		WhaleCount:   0,
		PriceDrop:    0,
		ATRCategory:  "NORMAL",
	}

	score := CalculateInterestScore(input)
	expected := 5.0
	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}
}

func TestCalculateInterestScore_PriceDropBoost(t *testing.T) {
	input := ScoringInput{
		CurrentPrice: 80,
		VWAPPrice:    100,
		ATRValue:     1.0,
		RSIValue:     50,
		WhaleCount:   0,
		PriceDrop:    20,
		ATRCategory:  "NORMAL",
	}
	score := CalculateInterestScore(input)
	expectedMin := 5.0 * 1.2
	if score < expectedMin {
		t.Errorf("Expected score at least %f, got %f", expectedMin, score)
	}
}

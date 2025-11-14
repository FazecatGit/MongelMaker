package strategy

import (
	"testing"

	"github.com/fazecat/mongelmaker/Internal/utils"
)

func TestCalculateInterestScore_NeutralConditions(t *testing.T) {
	input := utils.ScoringInput{
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
	input := utils.ScoringInput{
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

func TestCalculateInterestScore_HighWhaleActivity(t *testing.T) {
	input := utils.ScoringInput{
		CurrentPrice: 100,
		VWAPPrice:    100,
		ATRValue:     1.0,
		RSIValue:     50,
		WhaleCount:   10,
		PriceDrop:    0,
		ATRCategory:  "NORMAL",
	}
	score := CalculateInterestScore(input)
	expectedMin := 5.0 * 1.2
	if score < expectedMin {
		t.Errorf("Expected score at least %f, got %f", expectedMin, score)
	}
}

func TestCalculateInterestScore_LowRSIPenalty(t *testing.T) {
	input := utils.ScoringInput{
		CurrentPrice: 100,
		VWAPPrice:    100,
		ATRValue:     1.0,
		RSIValue:     25,
		WhaleCount:   0,
		PriceDrop:    0,
		ATRCategory:  "NORMAL",
	}
	score := CalculateInterestScore(input)
	expectedMax := 5.0 * 1.1
	if score > expectedMax {
		t.Errorf("Expected score at most %f, got %f", expectedMax, score)
	}
}

func TestCalculateInterestScore_HighATRAdjustment(t *testing.T) {
	input := utils.ScoringInput{
		CurrentPrice: 100,
		VWAPPrice:    100,
		ATRValue:     3.0,
		RSIValue:     50,
		WhaleCount:   0,
		PriceDrop:    0,
		ATRCategory:  "HIGH",
	}
	score := CalculateInterestScore(input)
	expectedMin := 5.0 * 1
	if score < expectedMin {
		t.Errorf("Expected score at least %f, got %f", expectedMin, score)
	}
}

package strategy

import "github.com/fazecat/mongelmaker/Internal/utils"

type ScoringInput struct {
	CurrentPrice float64
	VWAPPrice    float64
	ATRValue     float64
	RSIValue     float64
	WhaleCount   float64
	PriceDrop    float64
	ATRCategory  string
}

func CalculateInterestScore(input ScoringInput) float64 {
	baseScore := 5.0

	priceDropMult := 1.0 + (input.PriceDrop/10)*0.1

	vWapDistanceScore := (input.VWAPPrice - input.CurrentPrice) / input.VWAPPrice * 100
	var vwapMult float64
	if vWapDistanceScore >= 0 {
		vwapMult = 1.0
	} else {
		distanceAbs := utils.Abs(vWapDistanceScore)
		if distanceAbs <= 5 {
			vwapMult = 1.0 + (distanceAbs/5)*0.05
		} else if distanceAbs <= 15 {
			vwapMult = 1.05 + ((distanceAbs-5)/10)*0.10
		} else if distanceAbs <= 30 {
			vwapMult = 1.15 + ((distanceAbs-15)/15)*0.15
		} else {
			vwapMult = 1.3
		}
	}
	var atrMult float64
	if input.ATRCategory == "LOW" {
		atrMult = 1.2
	} else {
		atrMult = 1.0
	}

	whaleMult := 1.0
	if input.WhaleCount > 0 {
		whaleMult = 1.2
	}

	rsiMult := 1.0
	if input.RSIValue < 30 {
		rsiMult = 1.1
	}

	finalScore := baseScore * priceDropMult * vwapMult * atrMult * whaleMult * rsiMult
	return finalScore
}

package newsscraping

import (
	"regexp"
	"strings"
)

type CatalystDetector struct {
	earningsPatterns   []*regexp.Regexp
	acquisitionPattern []*regexp.Regexp
	regulatoryPattern  []*regexp.Regexp
	marketPattern      []*regexp.Regexp
}

func NewCatalystDetector() *CatalystDetector {
	return &CatalystDetector{
		earningsPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(earnings|q[1-4]\s*earnings|beat|miss|guidance)`),
			regexp.MustCompile(`(?i)(EPS|revenue|profit|margin)`),
		},
		acquisitionPattern: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(acquisition|acquired|merger|deal|buyout|takeover)`),
			regexp.MustCompile(`(?i)(partnership|joint\s*venture|collaboration)`),
		},
		regulatoryPattern: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(FDA|approval|clearance|ban|sanction|investigation)`),
			regexp.MustCompile(`(?i)(lawsuit|settlement|fine|regulation)`),
		},
		marketPattern: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(stock\s*split|dividend|buyback|delisting)`),
			regexp.MustCompile(`(?i)(IPO|spin-?off|offering)`),
		},
	}
}

func (cd *CatalystDetector) Detect(headline string) CatalystType {
	headlineLower := strings.ToLower(headline)

	for _, pattern := range cd.earningsPatterns {
		if pattern.MatchString(headlineLower) {
			return Earnings
		}
	}

	for _, pattern := range cd.acquisitionPattern {
		if pattern.MatchString(headlineLower) {
			return Acquisition
		}
	}

	for _, pattern := range cd.regulatoryPattern {
		if pattern.MatchString(headlineLower) {
			return Regulatory
		}
	}

	for _, pattern := range cd.marketPattern {
		if pattern.MatchString(headlineLower) {
			return Market
		}
	}

	return NoCatalyst
}

func (cd *CatalystDetector) GetImpact(catalystType CatalystType) float64 {
	impactMap := map[CatalystType]float64{
		Earnings:    0.15,
		Acquisition: 0.20,
		Regulatory:  0.25,
		Market:      0.10,
		Technical:   0.05,
		NoCatalyst:  0.02,
	}
	return impactMap[catalystType]
}

package data

import "strings"

// ModelPricing holds per-token prices in dollars per million tokens.
type ModelPricing struct {
	InputPerMTok       float64
	OutputPerMTok      float64
	CacheReadPerMTok   float64
	CacheWritePerMTok  float64
}

// pricingTable maps model family prefixes to pricing.
// Prices are in USD per million tokens.
var pricingTable = []struct {
	prefix  string
	pricing ModelPricing
}{
	{"claude-opus-4-6", ModelPricing{InputPerMTok: 15, OutputPerMTok: 75, CacheReadPerMTok: 1.50, CacheWritePerMTok: 18.75}},
	{"claude-opus-4-5", ModelPricing{InputPerMTok: 15, OutputPerMTok: 75, CacheReadPerMTok: 1.50, CacheWritePerMTok: 18.75}},
	{"claude-sonnet-4-5", ModelPricing{InputPerMTok: 3, OutputPerMTok: 15, CacheReadPerMTok: 0.30, CacheWritePerMTok: 3.75}},
	{"claude-sonnet-4-6", ModelPricing{InputPerMTok: 3, OutputPerMTok: 15, CacheReadPerMTok: 0.30, CacheWritePerMTok: 3.75}},
	{"claude-sonnet-4-", ModelPricing{InputPerMTok: 3, OutputPerMTok: 15, CacheReadPerMTok: 0.30, CacheWritePerMTok: 3.75}},
	{"claude-3-5-sonnet", ModelPricing{InputPerMTok: 3, OutputPerMTok: 15, CacheReadPerMTok: 0.30, CacheWritePerMTok: 3.75}},
	{"claude-3-5-haiku", ModelPricing{InputPerMTok: 0.80, OutputPerMTok: 4, CacheReadPerMTok: 0.08, CacheWritePerMTok: 1}},
	{"claude-haiku-4-5", ModelPricing{InputPerMTok: 0.80, OutputPerMTok: 4, CacheReadPerMTok: 0.08, CacheWritePerMTok: 1}},
	{"claude-3-opus", ModelPricing{InputPerMTok: 15, OutputPerMTok: 75, CacheReadPerMTok: 1.50, CacheWritePerMTok: 18.75}},
	// Fallback for unknown claude models — use Sonnet pricing as conservative middle ground
	{"claude-", ModelPricing{InputPerMTok: 3, OutputPerMTok: 15, CacheReadPerMTok: 0.30, CacheWritePerMTok: 3.75}},
}

// LookupPricing returns the pricing for a model name.
func LookupPricing(model string) ModelPricing {
	lower := strings.ToLower(model)
	for _, entry := range pricingTable {
		if strings.HasPrefix(lower, entry.prefix) {
			return entry.pricing
		}
	}
	// Default to Sonnet pricing if completely unknown
	return ModelPricing{InputPerMTok: 3, OutputPerMTok: 15, CacheReadPerMTok: 0.30, CacheWritePerMTok: 3.75}
}

// CalculateCost returns the estimated USD cost for a given token usage and model.
func CalculateCost(model string, usage ModelUsage) float64 {
	p := LookupPricing(model)
	cost := float64(usage.InputTokens) * p.InputPerMTok / 1_000_000
	cost += float64(usage.OutputTokens) * p.OutputPerMTok / 1_000_000
	cost += float64(usage.CacheReadInputTokens) * p.CacheReadPerMTok / 1_000_000
	cost += float64(usage.CacheCreationInputTokens) * p.CacheWritePerMTok / 1_000_000
	return cost
}

// TotalCost computes the total estimated cost across all models in a ModelUsage map.
func TotalCost(usage map[string]ModelUsage) float64 {
	var total float64
	for model, u := range usage {
		total += CalculateCost(model, u)
	}
	return total
}

package models

type TransformedCoin struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	MarketCap    float64 `json:"market_cap"`
	TotalVolume  float64 `json:"total_volume"`
	LastUpdated  string  `json:"last_updated"` //  ISO 8601 from API
}

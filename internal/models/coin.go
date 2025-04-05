package models

type CoinMarket struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	Image                    string  `json:"image"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCap                float64 `json:"market_cap"`
	MarketCapRank            int     `json:"market_cap_rank"`
	TotalVolume              float64 `json:"total_volume"`
	High24h                  float64 `json:"high_24h"`
	Low24h                   float64 `json:"low_24h"`
	PriceChange24h           float64 `json:"price_change_24h"`
	PriceChangePercent24h    float64 `json:"price_change_percentage_24h"`
	LastUpdated              string  `json:"last_updated"`
}

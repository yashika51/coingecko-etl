package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"coingecko-etl/internal/models"
)

func FetchMarketData() ([]models.CoinMarket, error) {
	baseURL := os.Getenv("COINGECKO_API_URL")
	vsCurrency := os.Getenv("VS_CURRENCY")
	perPage := os.Getenv("PER_PAGE")

	url := fmt.Sprintf("%s?vs_currency=%s&order=market_cap_desc&per_page=%s&page=1&sparkline=false",
		baseURL, vsCurrency, perPage)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("ERROR: Failed to fetch data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body: %v", err)
		return nil, err
	}

	var coins []models.CoinMarket
	err = json.Unmarshal(body, &coins)
	if err != nil {
		log.Printf("ERROR: Failed to unmarshal JSON: %v", err)
		return nil, err
	}

	log.Printf("INFO: Successfully fetched %d coin records", len(coins))
	return coins, nil
}

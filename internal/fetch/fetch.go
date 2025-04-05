package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"coingecko-etl/internal/models"
)

// FetchMarketData fetches coin market data from the CoinGecko API
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

// SaveRawData writes the coin data to a timestamped JSON file under data/raw/
func SaveRawData(coins []models.CoinMarket) error {
	err := os.MkdirAll("data/raw", os.ModePerm)
	if err != nil {
		log.Printf("ERROR: Failed to create raw data directory: %v", err)
		return err
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("data/raw/market_%s.json", timestamp)

	data, err := json.MarshalIndent(coins, "", "  ")
	if err != nil {
		log.Printf("ERROR: Failed to marshal raw data: %v", err)
		return err
	}

	err = os.WriteFile(filepath.Clean(filename), data, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write raw data file: %v", err)
		return err
	}

	log.Printf("INFO: Saved raw data to %s", filename)
	return nil
}

// filter and reshape the raw coin data
func TransformCoins(rawCoins []models.CoinMarket) []models.TransformedCoin {
	var processed []models.TransformedCoin

	for _, coin := range rawCoins {
		processed = append(processed, models.TransformedCoin{
			ID:           coin.ID,
			Symbol:       coin.Symbol,
			Name:         coin.Name,
			CurrentPrice: coin.CurrentPrice,
			MarketCap:    coin.MarketCap,
			TotalVolume:  coin.TotalVolume,
			LastUpdated:  coin.LastUpdated,
		})
	}
	return processed
}

func SaveProcessedData(coins []models.TransformedCoin) error {
	err := os.MkdirAll("data/processed", os.ModePerm)
	if err != nil {
		log.Printf("ERROR: Failed to create processed data directory: %v", err)
		return err
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("data/processed/market_%s.json", timestamp)

	data, err := json.MarshalIndent(coins, "", "  ")
	if err != nil {
		log.Printf("ERROR: Failed to marshal processed data: %v", err)
		return err
	}

	err = os.WriteFile(filepath.Clean(filename), data, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write processed data file: %v", err)
		return err
	}

	log.Printf("INFO: Saved processed data to %s", filename)
	return nil
}
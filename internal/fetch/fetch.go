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
	"coingecko-etl/internal/monitoring"
)

func FetchMarketData() ([]models.CoinMarket, error) {
	url := os.Getenv("COINGECKO_URL")
	if url == "" {
		url = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("ERROR: Failed to fetch data from API: %v", err)
		monitoring.FetchFailure.Inc()
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: API responded with status: %s", resp.Status)
		monitoring.FetchFailure.Inc()
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read API response: %v", err)
		monitoring.FetchFailure.Inc()
		return nil, err
	}

	var coins []models.CoinMarket
	if err := json.Unmarshal(body, &coins); err != nil {
		log.Printf("ERROR: Failed to unmarshal API response: %v", err)
		monitoring.FetchFailure.Inc()
		return nil, err
	}

	log.Printf("INFO: Successfully fetched %d coin records", len(coins))
	monitoring.FetchSuccess.Inc()
	monitoring.RecordsProcessed.Add(float64(len(coins)))
	return coins, nil
}

func SaveRawData(data []models.CoinMarket) error {
	if err := os.MkdirAll("data/raw", os.ModePerm); err != nil {
		return err
	}
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := filepath.Join("data/raw", fmt.Sprintf("market_%s.json", timestamp))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	log.Printf("INFO: Saved raw data to %s", filename)
	return nil
}

func SaveProcessedData(data []models.TransformedCoin) error {
	if err := os.MkdirAll("data/processed", os.ModePerm); err != nil {
		return err
	}
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := filepath.Join("data/processed", fmt.Sprintf("market_%s.json", timestamp))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	log.Printf("INFO: Saved processed data to %s", filename)
	return nil
}

func TransformCoins(raw []models.CoinMarket) []models.TransformedCoin {
	transformed := make([]models.TransformedCoin, 0, len(raw))
	for _, coin := range raw {
		transformed = append(transformed, models.TransformedCoin{
			ID:           coin.ID,
			Symbol:       coin.Symbol,
			Name:         coin.Name,
			CurrentPrice: coin.CurrentPrice,
			MarketCap:    coin.MarketCap,
			TotalVolume:  coin.TotalVolume,
			LastUpdated:  coin.LastUpdated,
		})
	}
	return transformed
}

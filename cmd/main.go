package main

import (
	"coingecko-etl/internal/fetch"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	coins, err := fetch.FetchMarketData()
	if err != nil {
		log.Fatalf("Failed to fetch market data: %v", err)
	}

	err = fetch.SaveRawData(coins)
	if err != nil {
		log.Fatalf("Failed to save raw data: %v", err)
	}

	transformed := fetch.TransformCoins(coins)

	err = fetch.SaveProcessedData(transformed)
	if err != nil {
		log.Fatalf("Failed to save processed data: %v", err)
	}

	log.Printf("First processed coin: %+v", transformed[0])
}

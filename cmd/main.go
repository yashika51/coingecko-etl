package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"coingecko-etl/internal/fetch"
)

func runETL() {
	coins, err := fetch.FetchMarketData()
	if err != nil {
		log.Printf("Failed to fetch market data: %v", err)
		return
	}

	if err := fetch.SaveRawData(coins); err != nil {
		log.Printf("Failed to save raw data: %v", err)
	}

	transformed := fetch.TransformCoins(coins)

	if err := fetch.SaveProcessedData(transformed); err != nil {
		log.Printf("Failed to save processed data: %v", err)
	}

	if err := fetch.SaveToPostgres(coins); err != nil {
		log.Printf("Failed to save to Postgres: %v", err)
	}

	log.Printf("ETL completed for %d coins", len(transformed))
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	runETL()

	for {
		select {
		case <-ticker.C:
			runETL()
		case <-ctx.Done():
			log.Println("Shutdown signal received. Exiting ETL loop.")
			return
		}
	}
}

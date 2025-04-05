package main

import (
	"context"
	"coingecko-etl/internal/fetch"
	"coingecko-etl/internal/utils"
	"coingecko-etl/internal/monitoring"
	"github.com/joho/godotenv"
	"log"
	"os/signal"
	"syscall"
	"time"
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
	utils.InitLogger()
	monitoring.Init()
	StartMetricsServer()


	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

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

package fetch

import (
	"context"
	"fmt"
	"log"
	"os"

	"coingecko-etl/internal/models"
	"github.com/jackc/pgx/v5"
)

func SaveToPostgres(coins []models.CoinMarket) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("ERROR: Postgres connection failed: %v", err)
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS coin_market_raw (
            id TEXT,
            symbol TEXT,
            name TEXT,
            image TEXT,
            current_price NUMERIC,
            market_cap NUMERIC,
            market_cap_rank INT,
            total_volume NUMERIC,
            high_24h NUMERIC,
            low_24h NUMERIC,
            price_change_24h NUMERIC,
            price_change_percentage_24h NUMERIC,
            last_updated TIMESTAMPTZ
        );
    `)
	if err != nil {
		log.Printf("ERROR: Failed to create table: %v", err)
		return err
	}

	for _, coin := range coins {
		_, err := conn.Exec(context.Background(), `
            INSERT INTO coin_market_raw (
                id, symbol, name, image, current_price, market_cap, market_cap_rank,
                total_volume, high_24h, low_24h, price_change_24h,
                price_change_percentage_24h, last_updated
            ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
        `,
			coin.ID, coin.Symbol, coin.Name, coin.Image,
			coin.CurrentPrice, coin.MarketCap, coin.MarketCapRank,
			coin.TotalVolume, coin.High24h, coin.Low24h,
			coin.PriceChange24h, coin.PriceChangePercent24h,
			coin.LastUpdated,
		)

		if err != nil {
			log.Printf("ERROR: Insert failed for %s: %v", coin.ID, err)
		}
	}

	var rowCount int
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM coin_market_raw").Scan(&rowCount)
	if err != nil {
		log.Printf("ERROR: Failed to validate row count: %v", err)
	} else {
		log.Printf("INFO: Total records in Postgres: %d", rowCount)
	}

	log.Printf("INFO: Inserted %d records into Postgres", len(coins))
	return nil
}

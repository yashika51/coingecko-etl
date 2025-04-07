# CoinGecko ETL Pipeline

This project implements an ETL (Extract, Transform, Load) pipeline in Go for ingesting cryptocurrency market data from the [CoinGecko API](https://www.coingecko.com/en/api). The project uses the free tier of CoinGecko’s API, so no authentication is required.

---

## Features

- Fetches cryptocurrency market data from CoinGecko every 30 seconds
- Saves the raw API response to timestamped `.json` files
- Transforms and normalizes the data structure for downstream use
- Saves both raw and processed data locally and to PostgreSQL
- Adds basic observability: structured logs, health checks, and Prometheus metrics
- Fully containerized with Docker and Docker Compose for local development


## Tech Stack

- **Language:** Go (Golang)
- **Data Source:** CoinGecko Open API
- **Database:** PostgreSQL (Dockerized)
- **Logging:** Standard Go `log` package
- **Monitoring:** Prometheus-compatible metrics exposed via `/metrics`
- **Scheduling:** Go `time.Ticker` for 30s intervals
- **Containerization:** Docker + Docker Compose

---

## Project Structure

```
.
├── cmd/                 # Entry point
├── internal/
│   ├── fetch/           # API fetch, transform, save, DB load
│   ├── models/          # Data structures for raw and processed JSON
│   ├── utils/           # Logger setup
│   └── monitoring/      # Prometheus metrics
├── data/                # Saved JSON files
├── logs/                # Logs
├── Dockerfile
├── docker-compose.yml
└── .env
```

---

## Setup Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/yashika51/coingecko-etl.git
   cd coingecko-etl
   ```

2. Create a `.env` file:
   ```env
   COINGECKO_URL=https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd
   VS_CURRENCY=usd
   PER_PAGE=10
   POSTGRES_HOST=db
   POSTGRES_PORT=5432
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=coingecko
   ```

3. Install Go dependencies:


Go modules will automatically install dependencies when you run `go build` or `go run`.
If needed, you can also install manually, for example:

   ```bash
   go get github.com/prometheus/client_golang/prometheus
   go get github.com/prometheus/client_golang/prometheus/promhttp
   ```

4. Run the system:
   ```bash
   docker compose up --build
   ```

---

## How to Verify

- Visit [http://localhost:8080/healthz](http://localhost:8080/healthz) → should return `OK`
- Visit [http://localhost:8080/metrics](http://localhost:8080/metrics) → Prometheus metrics
- Check `data/raw/` and `data/processed/` for timestamped `.json` files
- Log output goes to `stdout` or `logs/etl.log`
- Run `psql` inside the container to confirm data in `coin_market_raw`

---

## Sample Output

**Raw JSON:**

```json
[
  {
    "id": "bitcoin",
    "symbol": "btc",
    "name": "Bitcoin",
    "current_price": 70123,
    ...
  }
]
```

**Transformed JSON:**

```json
[
  {
    "id": "bitcoin",
    "symbol": "btc",
    "name": "Bitcoin",
    "current_price": 70123,
    "market_cap": 123456789,
    "total_volume": 8900000,
    "last_updated": "2025-04-05T14:34:01Z"
  }
]
```

---
## Productionization Plan


Here’s how I’d take this ETL system to production — making it scalable, observable, and easier to maintain long-term.

---

### Deployment Strategy

In a production setup, I would:

- Containerize the Go ETL service using Docker, and deploy it to a Kubernetes cluster (e.g. AKS or EKS).
- Schedule ETL runs using Kubernetes CronJobs or event-driven triggers like Azure Functions or AWS EventBridge.
- Use a managed PostgreSQL service like Amazon RDS or Azure Flexible Server, with automated backups and high availability across zones
- Expose Prometheus metrics from the service and visualize them using Grafana for monitoring system health, data freshness, and throughput.
- Centralize logs using cloud-native tools like AWS CloudWatch, Azure Monitor, with retention and search.
- Secure secrets using services like AWS Secrets Manager or Azure Key Vault.
- Automate builds and deployments using GitHub Actions or GitLab CI, pushing Docker images and deploying with Helm.

---

### Scalability

- I would modularize the pipeline into extract, transform, and load stages to allow independent scaling, parallelization, and reuse.
- If ingestion needs increase, a lightweight message queue like Kafka or GCP Pub/Sub could help decouple ingestion from processing.
- For long-term storage and analytics, I’d offload transformed data to columnar formats like Parquet in S3 or Azure Data Lake.
- Partitioning the Postgres table by date (or coin ID) could help optimize query performance and write throughput
- I’d consider implementing schema evolution safeguards and track transformation logic versions (e.g., in Git or dbt).

---

### Reliability & Fault Tolerance

- Implement retry logic with exponential backoff for API failures.
- Wrap database writes in transactions to avoid partial state on failure.
- Track metadata about each ETL run (timestamp, status, record count) to support debugging and monitoring.
- Handle shutdowns to avoid data loss during container termination or re-deployments.

---

### Cost Optimization

- Optimize polling frequency (e.g., reduce from every 30s to hourly during off-peak hours) to save on API calls and processing time.
- Compress raw and processed JSON using gzip or convert to columnar formats for archival.
- Archive historical data out of the primary database into cold storage or data lakes.
- Use cloud-native scheduling tools and serverless execution models to scale to zero when idle.

---

## Future Extensions

- I would like to add unit tests or integration tests to validate end-to-end ETL behavior
- Can also build a lightweight dashboard to visualize top coins or track ETL health
- Add some alerting on ETL failures like via email or Slack using Prometheus Alertmanager


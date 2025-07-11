# Marketflow

## Overview

Marketflow is a high-performance market data aggregator and processor designed for real-time streaming, aggregation, and storage of cryptocurrency or financial market data across multiple exchanges.

---

## Features

* 🔄 **Real-time data ingestion** from multiple exchanges
* ⏱ **Flexible aggregation** with arbitrary periods (not limited to minutes)
* 💾 **PostgreSQL storage** with optimized schema and indices
* ⚡ **Fast flushing and batching mechanisms**
* 🕒 **Timezone-safe period bucketing**
* 🛠 **Configurable pipeline** with easy period adjustments
* 📈 **Designed for scaling and integration** with analytical services

---

## Installation

```bash
# Clone repository
git clone https://platform.alem.school/azhaxyly/marketflow.git
cd marketflow

# Run
docker-compose up --build
```

### API Endpoints

**Market Data API**

`GET /prices/latest/{symbol}` – Get the latest price for a given symbol.

`GET /prices/latest/{exchange}/{symbol}` – Get the latest price for a given symbol from a specific exchange.

`GET /prices/highest/{symbol}` – Get the highest price over a period.

`GET /prices/highest/{exchange}/{symbol}` – Get the highest price over a period from a specific exchange.

`GET /prices/highest/{symbol}?period={duration}` – Get the highest price within the last `{duration}` (e.g., the last `1s`,  `3s`, `5s`, `10s`, `30s`, `1m`, `3m`, `5m`).

`GET /prices/highest/{exchange}/{symbol}?period={duration}` – Get the highest price within the last `{duration}` from a specific exchange.

`GET /prices/lowest/{symbol}` – Get the lowest price over a period.

`GET /prices/lowest/{exchange}/{symbol}` – Get the lowest price over a period from a specific exchange.

`GET /prices/lowest/{symbol}?period={duration}` – Get the lowest price within the last {duration}.

`GET /prices/lowest/{exchange}/{symbol}?period={duration}` – Get the lowest price within the last `{duration}` from a specific exchange.

`GET /prices/average/{symbol}` – Get the average price over a period.

`GET /prices/average/{exchange}/{symbol}` – Get the average price over a period from a specific exchange.

`GET /prices/average/{exchange}/{symbol}?period={duration}` – Get the average price within the last `{duration}` from a specific exchange

**Data Mode API**

`POST /mode/test` – Switch to `Test Mode` (use generated data).

`POST /mode/live` – Switch to `Live Mode` (fetch data from `provided programs`).

**System Health**

`GET /health` - Returns system status (e.g., connections, Redis availability).  

## Authors

MarketFlow is maintained by **azhaxyly** and **mromanul**. Contributions are welcome via pull requests.

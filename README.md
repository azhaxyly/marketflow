# MarketFlow

Real-Time Market Data Processing System for cryptocurrency prices.

## Setup

1. Copy the example config file: `cp configs/config.example.yaml configs/config.yaml`
2. Customize the configuration as needed
3. Start dependencies: `docker-compose up -d`
4. Build the application: `go build -o marketflow cmd/marketflow/main.go`
5. Run the application: `./marketflow`

## Test Mode Generator

To build and run the test data generator:

```
go build -o testgenerator cmd/testgenerator/main.go
./testgenerator
```

## API Endpoints

See the project documentation for available API endpoints.

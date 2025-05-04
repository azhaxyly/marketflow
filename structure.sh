#!/bin/bash

echo "Initializing marketflow project structure..."

# Create main directories
mkdir -p cmd/marketflow
mkdir -p cmd/testgenerator
mkdir -p internal/domain/model
mkdir -p internal/app/usecase
mkdir -p internal/adapter/web
mkdir -p internal/adapter/storage/postgres
mkdir -p internal/adapter/cache/redis
mkdir -p internal/adapter/exchange
mkdir -p internal/port
mkdir -p configs
mkdir -p scripts
mkdir -p test

# Create main Go files
touch cmd/marketflow/main.go
touch cmd/testgenerator/main.go

# Create config files
cat > configs/config.example.yaml << EOF
server:
  port: 8080

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: marketflow
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

exchanges:
  exchange1:
    host: localhost
    port: 40101
  exchange2:
    host: localhost
    port: 40102
  exchange3:
    host: localhost
    port: 40103

testmode:
  host: localhost
  port: 40104
EOF

touch configs/config.yaml

# Create domain models
touch internal/domain/model/price.go
touch internal/domain/model/exchange.go

# Create use cases
touch internal/app/usecase/price_service.go
touch internal/app/usecase/exchange_service.go
touch internal/app/usecase/mode_service.go

# Create adapters
touch internal/adapter/web/handler.go
touch internal/adapter/web/router.go
touch internal/adapter/storage/postgres/price_repository.go
touch internal/adapter/cache/redis/price_cache.go
touch internal/adapter/exchange/client.go
touch internal/adapter/exchange/listener.go

# Create ports
touch internal/port/repository.go
touch internal/port/cache.go
touch internal/port/exchange.go

# Create Docker and Docker Compose files
cat > docker-compose.yml << EOF
version: '3'

services:
  postgres:
    image: postgres:latest
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: marketflow
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:latest
    ports:
      - "6379:6379"

volumes:
  postgres_data:
EOF

# Create SQL initialization script
mkdir -p scripts
cat > scripts/init.sql << EOF
CREATE TABLE IF NOT EXISTS price_aggregates (
    id SERIAL PRIMARY KEY,
    pair_name VARCHAR(20) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    average_price DECIMAL(24,8) NOT NULL,
    min_price DECIMAL(24,8) NOT NULL,
    max_price DECIMAL(24,8) NOT NULL
);

CREATE INDEX idx_pair_timestamp ON price_aggregates(pair_name, timestamp);
CREATE INDEX idx_exchange_pair_timestamp ON price_aggregates(exchange, pair_name, timestamp);
EOF

# Create Go module
go mod init marketflow
go mod tidy

# Create a basic README
cat > README.md << EOF
# MarketFlow

Real-Time Market Data Processing System for cryptocurrency prices.

## Setup

1. Copy the example config file: \`cp configs/config.example.yaml configs/config.yaml\`
2. Customize the configuration as needed
3. Start dependencies: \`docker-compose up -d\`
4. Build the application: \`go build -o marketflow cmd/marketflow/main.go\`
5. Run the application: \`./marketflow\`

## Test Mode Generator

To build and run the test data generator:

\`\`\`
go build -o testgenerator cmd/testgenerator/main.go
./testgenerator
\`\`\`

## API Endpoints

See the project documentation for available API endpoints.
EOF

# Create .gitignore
cat > .gitignore << EOF
# Binary files
marketflow
testgenerator

# Config file with potentially sensitive information
configs/config.yaml

# IDE files
.idea/
.vscode/

# Go files
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work

# Dependency directories
vendor/
EOF

echo "Project structure initialized successfully!"
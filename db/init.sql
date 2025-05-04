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

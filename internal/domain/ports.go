package domain

import (
	"context"
	"time"
)

// redis
type Cache interface {
	SetLatest(update PriceUpdate) error
	GetLatest(pair string) (PriceUpdate, error)
}

// postgres
type PriceRepository interface {
    StoreStats(stat PriceStats) error
    StoreStatsBatch(stats []PriceStats) error
    GetStats(pair, exchange string, since time.Time) ([]PriceStats, error)
    GetLatest(ctx context.Context, exchange, pair string) (PriceStats, error)
    GetByPeriod(ctx context.Context, exchange, pair string, period time.Duration) ([]PriceStats, error)
}

// http
type ExchangeClient interface {
	Start(ctx context.Context, out chan<- PriceUpdate) error
	Stop() error
}

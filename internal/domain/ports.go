package domain

import "time"

// redis
type Cache interface {
	SetLatest(update PriceUpdate) error
	GetLatest(pair string) (PriceUpdate, error)
}

// postgres
type PriceRepository interface {
	StoreStats(stat PriceStats) error
	GetStats(pair, exchange string, since time.Time) ([]PriceStats, error)
}

// http
type ExchangeClient interface {
	Start(out chan<- PriceUpdate) error
	Stop() error
}

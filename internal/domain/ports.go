package domain

import "time"

type Cache interface {
	SetLatest(update PriceUpdate) error
	GetLatest(pair string) (PriceUpdate, error)
}

type PriceRepository interface {
	StoreStats(stat PriceStats) error
	GetStats(pair, exchange string, since time.Time) ([]PriceStats, error)
}

type ExchangeClient interface {
	Start(out chan<- PriceUpdate) error
	Stop() error
}

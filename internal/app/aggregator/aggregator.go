package aggregator

import (
	"context"
	"time"

	"marketflow/internal/adapters/redis"
	"marketflow/internal/domain"
	"marketflow/internal/logger"
)

type Aggregator struct {
	Input  <-chan domain.PriceUpdate
	Repo   domain.PriceRepository
	Cache  *redis.RedisCache
	Window time.Duration
}

func NewAggregator(input <-chan domain.PriceUpdate, repo domain.PriceRepository, cache *redis.RedisCache, window time.Duration) *Aggregator {
	return &Aggregator{
		Input:  input,
		Repo:   repo,
		Cache:  cache,
		Window: window,
	}
}

func (a *Aggregator) Start(ctx context.Context) {
	buffer := make(map[string][]float64)
	ticker := time.NewTicker(a.Window)
	cleanTicker := time.NewTicker(5 * time.Minute) // Очистка каждые 5 минут
	defer ticker.Stop()
	defer cleanTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.flush(buffer, time.Now())
			return
		case update, ok := <-a.Input:
			if !ok {
				a.flush(buffer, time.Now())
				return
			}
			key := update.Exchange + ":" + update.Pair
			buffer[key] = append(buffer[key], update.Price)

		case tickTime := <-ticker.C:
			a.flush(buffer, tickTime)
			buffer = make(map[string][]float64)

		case <-cleanTicker.C:
			if err := a.Cache.CleanOld(ctx, "latest:*"); err != nil {
				logger.Error("failed to clean old redis keys", "error", err)
			}
		}
	}
}

func (a *Aggregator) flush(buffer map[string][]float64, ts time.Time) {
	var stats []domain.PriceStats
	for key, prices := range buffer {
		if len(prices) == 0 {
			continue
		}
		parts := splitKey(key)
		exchange, pair := parts[0], parts[1]

		var sum, min, max float64
		min = prices[0]
		max = prices[0]
		for _, p := range prices {
			sum += p
			if p < min {
				min = p
			}
			if p > max {
				max = p
			}
		}
		avg := sum / float64(len(prices))

		stat := domain.PriceStats{
			Exchange:  exchange,
			Pair:      pair,
			Timestamp: ts,
			Average:   avg,
			Min:       min,
			Max:       max,
		}
		stats = append(stats, stat)
	}

	if len(stats) > 0 {
		if err := a.Repo.StoreStatsBatch(stats); err != nil {
			logger.Error("failed to store batch stats", "error", err)
		} else {
			logger.Info("stored batch stats", "count", len(stats))
		}
	}
}

func splitKey(key string) []string {
	for i := 0; i < len(key); i++ {
		if key[i] == ':' {
			return []string{key[:i], key[i+1:]}
		}
	}
	return []string{"", key}
}
package aggregator

import (
	"time"

	"marketflow/internal/domain"
)

type Aggregator struct {
	Input  <-chan domain.PriceUpdate
	Repo   domain.PriceRepository
	Window time.Duration
}

func NewAggregator(input <-chan domain.PriceUpdate, repo domain.PriceRepository, window time.Duration) *Aggregator {
	return &Aggregator{
		Input:  input,
		Repo:   repo,
		Window: window,
	}
}

func (a *Aggregator) Start() {
	buffer := make(map[string][]float64)
	ticker := time.NewTicker(a.Window)
	defer ticker.Stop()

	for {
		select {
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
		}
	}
}

func (a *Aggregator) flush(buffer map[string][]float64, ts time.Time) {
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
		err := a.Repo.StoreStats(stat)
		if err != nil {
			// log and continue
			// ideally use structured logger
			println("[aggregator] failed to store stat", err.Error())
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

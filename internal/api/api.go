package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"marketflow/internal/adapters/redis"
	"marketflow/internal/app/mode"
	"marketflow/internal/domain"
	"marketflow/internal/logger"
)

type Server struct {
	repo    domain.PriceRepository
	cache   *redis.RedisCache
	manager *mode.Manager
}

func NewServer(repo domain.PriceRepository, cache *redis.RedisCache, manager *mode.Manager) *Server {
	return &Server{
		repo:    repo,
		cache:   cache,
		manager: manager,
	}
}

func (s *Server) Router(input chan<- domain.PriceUpdate) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/prices/latest/", s.handleLatestPrice(input))
	mux.HandleFunc("/prices/highest/", s.handleHighestPrice)
	mux.HandleFunc("/prices/lowest/", s.handleLowestPrice)
	mux.HandleFunc("/prices/average/", s.handleAveragePrice)
	mux.HandleFunc("/mode/test", s.handleSetTestMode(input))
	mux.HandleFunc("/mode/live", s.handleSetLiveMode(input))
	mux.HandleFunc("/health", s.handleHealth)

	return mux
}

func (s *Server) handleLowestPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	symbol := parts[3]
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "ex1"
	}
	periodStr := r.URL.Query().Get("period")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	stats, err := s.repo.GetByPeriod(ctx, exchange, symbol, period)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	var minPrice float64
	var minTime time.Time
	for i, stat := range stats {
		if i == 0 || stat.Min < minPrice {
			minPrice = stat.Min
			minTime = stat.Timestamp
		}
	}

	if len(stats) == 0 {
		http.Error(w, "no data for period", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"exchange": exchange,
		"pair":     symbol,
		"price":    minPrice,
		"time":     minTime,
	})
}

func (s *Server) handleAveragePrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	symbol := parts[3]
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "ex1"
	}
	periodStr := r.URL.Query().Get("period")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	stats, err := s.repo.GetByPeriod(ctx, exchange, symbol, period)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	var sum float64
	for _, stat := range stats {
		sum += stat.Average
	}
	if len(stats) == 0 {
		http.Error(w, "no data for period", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"exchange": exchange,
		"pair":     symbol,
		"price":    sum / float64(len(stats)),
		"count":    len(stats),
	})
}

func (s *Server) handleLatestPrice(input chan<- domain.PriceUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		symbol := parts[3]
		exchange := r.URL.Query().Get("exchange")
		if exchange == "" {
			exchange = "ex1" // Default exchange
		}

		update, err := s.cache.GetLatest(ctx, exchange, symbol)
		if err != nil {
			logger.Warn("cache miss, falling back to postgres", "symbol", symbol, "exchange", exchange)
			stats, err := s.repo.GetLatest(ctx, exchange, symbol)
			if err != nil {
				logger.Error("failed to get latest price", "symbol", symbol, "exchange", exchange, "error", err)
				http.Error(w, "failed to get latest price", http.StatusInternalServerError)
				return
			}
			update = domain.PriceUpdate{
				Exchange: stats.Exchange,
				Pair:     stats.Pair,
				Price:    stats.Average,
				Time:     stats.Timestamp,
			}
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"exchange": update.Exchange,
			"pair":     update.Pair,
			"price":    update.Price,
			"time":     update.Time,
		})
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, errRedis := s.cache.GetLatest(ctx, "ex1", "BTCUSDT")
	_, errPg := s.repo.GetLatest(ctx, "ex1", "BTCUSDT")

	status := map[string]string{
		"redis":    "ok",
		"postgres": "ok",
	}
	if errRedis != nil {
		status["redis"] = "unavailable"
	}
	if errPg != nil {
		status["postgres"] = "unavailable"
	}

	respondJSON(w, http.StatusOK, status)
}

func (s *Server) handleHighestPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	symbol := parts[3]
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "ex1" // Default exchange
	}

	periodStr := r.URL.Query().Get("period")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		logger.Error("invalid period", "period", periodStr, "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	stats, err := s.repo.GetByPeriod(ctx, exchange, symbol, period)
	if err != nil {
		logger.Error("failed to get stats by period", "symbol", symbol, "exchange", exchange, "period", period, "error", err)
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	var maxPrice float64
	var maxTime time.Time
	for _, stat := range stats {
		if stat.Max > maxPrice {
			maxPrice = stat.Max
			maxTime = stat.Timestamp
		}
	}

	if maxPrice == 0 {
		http.Error(w, "no data for period", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"exchange": exchange,
		"pair":     symbol,
		"price":    maxPrice,
		"time":     maxTime,
	})
}

func (s *Server) handleSetTestMode(input chan<- domain.PriceUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := s.manager.Start(r.Context(), input, mode.Test); err != nil {
			logger.Error("failed to set test mode", "error", err)
			http.Error(w, "failed to set test mode", http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"mode": "test"})
	}
}

func (s *Server) handleSetLiveMode(input chan<- domain.PriceUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := s.manager.Start(r.Context(), input, mode.Live); err != nil {
			logger.Error("failed to set live mode", "error", err)
			http.Error(w, "failed to set live mode", http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"mode": "live"})
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

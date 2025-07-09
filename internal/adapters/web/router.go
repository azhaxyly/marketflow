package web

import (
	"marketflow/internal/domain"
	"net/http"
)

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

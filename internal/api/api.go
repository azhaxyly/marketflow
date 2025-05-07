package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"marketflow/internal/app/mode"
	"marketflow/internal/domain"
	"marketflow/internal/logger"
)

func HandleLiveMode(manager *mode.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := manager.Start(r.Context(), make(chan domain.PriceUpdate, 1000), mode.Live); err != nil {
			logger.Error("failed to switch to live mode", "error", err)
			http.Error(w, fmt.Sprintf("Error starting live mode: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Switched to Live mode"))
	}
}

func HandleTestMode(manager *mode.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := manager.Start(r.Context(), make(chan domain.PriceUpdate, 1000), mode.Test); err != nil {
			logger.Error("failed to switch to test mode", "error", err)
			http.Error(w, fmt.Sprintf("Error starting test mode: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Switched to Test mode"))
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь добавить проверку состояния Redis и других сервисов
	w.WriteHeader(http.StatusOK)
	logger.Info("health check passed")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

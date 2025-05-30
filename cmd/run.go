package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketflow/internal/adapters/redis"
	"marketflow/internal/adapters/storage/postgres"
	"marketflow/internal/api"
	"marketflow/internal/app/aggregator"
	"marketflow/internal/app/mode"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/logger"
)

func Run() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	logger.Init(env)

	logger.Info("starting application", "env", env)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pgDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DBName, cfg.Postgres.SSLMode)

	repo, err := postgres.NewPostgresRepository(pgDSN)
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}
	defer repo.Close()

	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	cache := redis.NewRedisCache(redisAddr, cfg.Redis.Password, cfg.Redis.DB, cfg.RedisTTL)
	defer cache.Close()

	inputChan := make(chan domain.PriceUpdate, 1000)
	manager := mode.NewManager(cfg)
	agg := aggregator.NewAggregator(inputChan, repo, cache, cfg.AggregatorWindow)

	go agg.Start(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := manager.Start(ctx, inputChan, mode.Test); err != nil {
		log.Fatalf("failed to start test mode: %v", err)
	}

	apiServer := api.NewServer(repo, cache, manager)

	srv := &http.Server{
		Addr:    cfg.APIAddr,
		Handler: apiServer.Router(inputChan), // создадим метод Router() чуть ниже
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("API server failed", "error", err)
			log.Fatalf("API server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("API shutdown error", "error", err)
	}

	cancel()

	logger.Info("shutdown complete")
}

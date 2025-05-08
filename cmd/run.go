package cmd

import (
	"context"
	"fmt"
	"log"

	"marketflow/internal/adapters/redis"
	"marketflow/internal/adapters/storage/postgres"
	"marketflow/internal/api"
	"marketflow/internal/app/aggregator"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/logger"
	"marketflow/internal/app/mode" 
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Формирование Postgres DSN
	pgDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DBName, cfg.Postgres.SSLMode)

	// Инициализация зависимостей
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

	// Запуск Aggregator
	go agg.Start(context.Background())

	// Запуск в Test Mode по умолчанию
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := manager.Start(ctx, inputChan, mode.Test); err != nil {
		log.Fatalf("failed to start test mode: %v", err)
	}

	// Запуск API
	apiServer := api.NewServer(repo, cache, manager)
	go func() {
		if err := apiServer.Start(cfg.APIAddr, inputChan); err != nil {
			logger.Error("API server failed", "error", err)
			log.Fatalf("API server failed: %v", err)
		}
	}()

	// Блокируем основной поток
	select {}
}
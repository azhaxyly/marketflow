package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"marketflow/internal/domain"
	"marketflow/internal/logger"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	cache := &RedisCache{
		client: client,
		ttl:    ttl,
	}

	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			logger.Error("failed to ping redis", "attempt", i+1, "error", err)
			if i == 2 {
				logger.Warn("redis connection failed after retries, proceeding with fallback")
			}
			time.Sleep(time.Second * time.Duration(i+1))
		} else {
			logger.Info("redis connection established")
			break
		}
	}

	return cache
}

func (r *RedisCache) SetLatest(ctx context.Context, update domain.PriceUpdate) error {
	key := fmt.Sprintf("latest:%s:%s", update.Exchange, update.Pair)
	data, err := json.Marshal(update)
	if err != nil {
		logger.Error("marshal error", "key", key, "error", err)
		return fmt.Errorf("marshal error: %w", err)
	}
	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		logger.Warn("redis set error, using fallback", "key", key, "error", err)
		return nil
	}
	logger.Info("updated latest price", "key", key, "price", update.Price)
	return nil
}

func (r *RedisCache) GetLatest(ctx context.Context, exchange, pair string) (domain.PriceUpdate, error) {
	key := fmt.Sprintf("latest:%s:%s", exchange, pair)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Warn("no data in redis", "key", key)
		return domain.PriceUpdate{}, fmt.Errorf("no data for %s", key)
	}
	if err != nil {
		logger.Warn("redis get error, using fallback", "key", key, "error", err)
		return domain.PriceUpdate{}, fmt.Errorf("redis unavailable: %w", err)
	}
	var update domain.PriceUpdate
	if err := json.Unmarshal([]byte(val), &update); err != nil {
		logger.Error("unmarshal error", "key", key, "error", err)
		return domain.PriceUpdate{}, fmt.Errorf("unmarshal error: %w", err)
	}
	logger.Info("got latest price", "key", key, "price", update.Price)
	return update, nil
}

func (r *RedisCache) CleanOld(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		logger.Warn("failed to scan keys for cleanup", "pattern", pattern, "error", err)
		return nil
	}
	if len(keys) == 0 {
		return nil
	}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		logger.Warn("failed to delete old keys", "pattern", pattern, "error", err)
		return nil
	}
	logger.Info("cleaned old keys", "pattern", pattern, "count", len(keys))
	return nil
}

func (r *RedisCache) Close() error {
	logger.Info("closing redis cache")
	return r.client.Close()
}

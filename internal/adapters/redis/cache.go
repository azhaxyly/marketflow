package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"marketflow/internal/domain"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{client: rdb, ttl: ttl}
}

func (r *RedisCache) SetLatest(update domain.PriceUpdate) error {
	ctx := context.Background()
	key := fmt.Sprintf("latest:%s:%s", update.Exchange, update.Pair)
	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal PriceUpdate: %w", err)
	}
	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set Redis key %s: %w", key, err)
	}
	return nil
}

func (r *RedisCache) GetLatest(pair string) (domain.PriceUpdate, error) {
	ctx := context.Background()
	key := pair
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return domain.PriceUpdate{}, fmt.Errorf("failed to get Redis key %s: %w", key, err)
	}

	var update domain.PriceUpdate
	if err := json.Unmarshal([]byte(val), &update); err != nil {
		return domain.PriceUpdate{}, fmt.Errorf("failed to unmarshal PriceUpdate: %w", err)
	}
	return update, nil
}

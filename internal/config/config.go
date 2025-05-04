package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Postgres  PostgresConfig
	Redis     RedisConfig
	Exchanges []Exchange
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type Exchange struct {
	Name    string
	Address string
}

func Load() (*Config, error) {
	pgPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_PORT: %w", err)
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}

	cfg := &Config{
		Postgres: PostgresConfig{
			Host:     os.Getenv("PG_HOST"),
			Port:     pgPort,
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   os.Getenv("PG_DB"),
			SSLMode:  os.Getenv("PG_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     redisPort,
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Exchanges: []Exchange{
			{Name: "exchange1", Address: os.Getenv("EXCHANGE1_ADDR")},
			{Name: "exchange2", Address: os.Getenv("EXCHANGE2_ADDR")},
			{Name: "exchange3", Address: os.Getenv("EXCHANGE3_ADDR")},
		},
	}

	return cfg, nil
}

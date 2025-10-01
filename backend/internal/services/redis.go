package services

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/macielcr7/map-area-factions/backend/internal/config"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(cfg *config.RedisConfig) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisService{Client: client}, nil
}

func (rs *RedisService) Close() error {
	return rs.Client.Close()
}

func (rs *RedisService) HealthCheck() error {
	ctx := context.Background()
	return rs.Client.Ping(ctx).Err()
}
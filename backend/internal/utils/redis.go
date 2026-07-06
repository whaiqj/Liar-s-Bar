package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"liars-bar/internal/config"
)

var Rdb *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}
	log.Println("Redis connected successfully")
	return nil
}

func CacheSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Rdb.Set(ctx, key, data, expiration).Err()
}

func CacheGet(ctx context.Context, key string, dest interface{}) error {
	data, err := Rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func CacheDel(ctx context.Context, keys ...string) error {
	return Rdb.Del(ctx, keys...).Err()
}

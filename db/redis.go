package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ritikchawla/url-shortner/config"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("error connecting to Redis: %v", err)
	}

	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

// Cache URL mapping
func CacheURL(shortCode, longURL string, expiration time.Duration) error {
	return RedisClient.Set(ctx, shortCode, longURL, expiration).Err()
}

// Get cached URL
func GetCachedURL(shortCode string) (string, error) {
	longURL, err := RedisClient.Get(ctx, shortCode).Result()
	if err == redis.Nil {
		return "", nil // URL not found in cache
	}
	return longURL, err
}

// Increment visit count in Redis
func IncrementVisits(shortCode string) error {
	return RedisClient.Incr(ctx, fmt.Sprintf("visits:%s", shortCode)).Err()
}

// Get visit count from Redis
func GetVisits(shortCode string) (int64, error) {
	visits, err := RedisClient.Get(ctx, fmt.Sprintf("visits:%s", shortCode)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return visits, err
}

package nosql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yoockh/dbyoc/logger"
	"github.com/yoockh/dbyoc/utils"
)

type RedisClient struct {
	client *redis.Client
	logger logger.Logger
}

func NewRedisClient(addr string, password string, db int, log logger.Logger) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		client: rdb,
		logger: log,
	}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set value in Redis", "key", key, "error", err)
		return err
	}
	return nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Warn("Key does not exist", "key", key)
			return "", nil
		}
		r.logger.Error("Failed to get value from Redis", "key", key, "error", err)
		return "", err
	}
	return val, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		r.logger.Error("Failed to ping Redis", "error", err)
		return err
	}
	return nil
}

func (r *RedisClient) RetryOperation(ctx context.Context, operation func() error) error {
	return utils.Retry(operation)
}

func (r *RedisClient) Reconnect() error {
	r.logger.Info("Reconnecting to Redis...")
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.client.Options().Addr,
		Password: r.client.Options().Password,
		DB:       r.client.Options().DB,
	})
	return nil
}

package nosql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/yoockh/dbyoc/config"
	"github.com/yoockh/dbyoc/utils"
)

type RedisClient struct {
	client *redis.Client
	logger *logrus.Logger
	config config.RedisConfig
}

func NewRedisClient(cfg config.RedisConfig, log *logrus.Logger) *RedisClient {
	var opts *redis.Options

	// Prioritize URL if provided
	if cfg.URL != "" {
		opts, _ = redis.ParseURL(cfg.URL)
	} else {
		opts = &redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}
	}

	rdb := redis.NewClient(opts)

	return &RedisClient{
		client: rdb,
		logger: log,
		config: cfg,
	}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		r.logger.WithFields(logrus.Fields{"key": key, "error": err}).Error("Failed to set value in Redis")
		return err
	}
	return nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.WithField("key", key).Warn("Key does not exist")
			return "", nil
		}
		r.logger.WithFields(logrus.Fields{"key": key, "error": err}).Error("Failed to get value from Redis")
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
		r.logger.WithField("error", err).Error("Failed to ping Redis")
		return err
	}
	return nil
}

func (r *RedisClient) RetryOperation(ctx context.Context, operation func() error) error {
	return utils.Retry(operation, nil)
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

// QuickRedis creates Redis client from REDIS_URL env only
func QuickRedis(logger ...*logrus.Logger) (*RedisClient, error) {
	cfg, err := config.QuickRedisConfig()
	if err != nil {
		return nil, err
	}

	log := logrus.New()
	if len(logger) > 0 {
		log = logger[0]
	}

	return NewRedisClient(*cfg, log), nil
}

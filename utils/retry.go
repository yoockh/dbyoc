package utils

import (
	"errors"
	"fmt"
	"time"
)

var ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")

type RetryConfig struct {
	MaxRetries int
	Delay      time.Duration
}

func NewRetryConfig(maxRetries int, delay time.Duration) *RetryConfig {
	return &RetryConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
	}
}

func Retry(operation func() error, config *RetryConfig) error {
	if config == nil {
		config = NewRetryConfig(3, 100*time.Millisecond)
	}

	var lastErr error
	for i := 0; i < config.MaxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
			time.Sleep(config.Delay)
		}
	}
	return fmt.Errorf("%w: %v", ErrMaxRetriesExceeded, lastErr)
}

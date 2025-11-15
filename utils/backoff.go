package utils

import (
	"math"
	"time"
)

// Backoff defines the structure for exponential backoff.
type Backoff struct {
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
}

// NewBackoff creates a new Backoff instance with default values.
func NewBackoff() *Backoff {
	return &Backoff{
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     30 * time.Second,
		MaxElapsedTime:  5 * time.Minute,
	}
}

// GetNextInterval calculates the next backoff interval.
func (b *Backoff) GetNextInterval(attempt int) time.Duration {
	interval := time.Duration(float64(b.InitialInterval) * math.Pow(b.Multiplier, float64(attempt)))
	if interval > b.MaxInterval {
		interval = b.MaxInterval
	}
	return interval
}

// IsElapsed checks if the maximum elapsed time has been reached.
func (b *Backoff) IsElapsed(start time.Time) bool {
	return time.Since(start) >= b.MaxElapsedTime
}
package metrics

import (
    "sync"
    "time"
)

type MetricsCollector struct {
    mu              sync.Mutex
    queryCount      int
    connectionCount int
    startTime       time.Time
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        startTime: time.Now(),
    }
}

func (mc *MetricsCollector) IncrementQueryCount() {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    mc.queryCount++
}

func (mc *MetricsCollector) IncrementConnectionCount() {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    mc.connectionCount++
}

func (mc *MetricsCollector) GetMetrics() (int, int, time.Duration) {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    uptime := time.Since(mc.startTime)
    return mc.queryCount, mc.connectionCount, uptime
}

func (mc *MetricsCollector) Reset() {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    mc.queryCount = 0
    mc.connectionCount = 0
    mc.startTime = time.Now()
}
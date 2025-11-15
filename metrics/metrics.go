package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu                sync.Mutex
	TotalConnections  int
	ActiveConnections int
	QueryCount        int
	TotalQueryTime    time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncrementConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalConnections++
	m.ActiveConnections++
}

func (m *Metrics) DecrementConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveConnections--
}

func (m *Metrics) RecordQuery(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.QueryCount++
	m.TotalQueryTime += duration
}

func (m *Metrics) GetMetrics() (int, int, int, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.TotalConnections, m.ActiveConnections, m.QueryCount, m.TotalQueryTime
}
package metrics

type QueryMetrics struct {
	QueryString   string
	ExecutionTime float64
	Success       bool
}

type ConnectionMetrics struct {
	ConnectionID string
	IsActive     bool
	LastUsed     int64
}

package metrics

type Metrics struct {
    TotalQueries     int64
    SuccessfulQueries int64
    FailedQueries    int64
    ConnectionCount  int64
}

type QueryMetrics struct {
    QueryString string
    ExecutionTime float64
    Success bool
}

type ConnectionMetrics struct {
    ConnectionID string
    IsActive     bool
    LastUsed     int64
}
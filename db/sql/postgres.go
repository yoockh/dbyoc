package sql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/yoockh/dbyoc/config"
)

type PostgresClient struct {
	*sql.DB
	dbConfig config.DatabaseConfig
}

func NewPostgresClient(cfg config.DatabaseConfig) (*PostgresClient, error) {
	var connStr string

	// Prioritize URL if provided
	if cfg.URL != "" {
		connStr = cfg.URL
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Set pool settings
	if cfg.MaxPoolSize > 0 {
		db.SetMaxOpenConns(cfg.MaxPoolSize)
		db.SetMaxIdleConns(cfg.MaxPoolSize / 2)
	}
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresClient{DB: db, dbConfig: cfg}, nil
}

func (p *PostgresClient) Close() error {
	return p.DB.Close()
}

func (p *PostgresClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.Query(query, args...)
}

func (p *PostgresClient) Insert(query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

func (p *PostgresClient) Update(query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

func (p *PostgresClient) RetryQuery(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	maxRetries := p.dbConfig.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	for i := 0; i < maxRetries; i++ {
		rows, err = p.Query(query, args...)
		if err == nil {
			return rows, nil
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	return nil, fmt.Errorf("query failed after %d attempts: %w", maxRetries, err)
}

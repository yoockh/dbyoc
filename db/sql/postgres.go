package sql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/yourusername/dbyoc/db"
)

type PostgresClient struct {
	*sql.DB
	config db.Config
}

func NewPostgresClient(config db.Config) (*PostgresClient, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresClient{DB: db, config: config}, nil
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
	for i := 0; i < p.config.MaxRetries; i++ {
		rows, err = p.Query(query, args...)
		if err == nil {
			return rows, nil
		}
		time.Sleep(time.Duration(i) * time.Second) // Exponential backoff can be implemented here
	}
	return nil, fmt.Errorf("query failed after %d attempts: %w", p.config.MaxRetries, err)
}
package db

import (
	"database/sql"
	"sync"
	"time"
)

type DBPool struct {
	pool       *sql.DB
	maxOpen    int
	maxIdle    int
	idleTimeout time.Duration
	mu         sync.Mutex
}

func NewDBPool(dataSourceName string, maxOpen, maxIdle int, idleTimeout time.Duration) (*DBPool, error) {
	db, err := sql.Open("postgres", dataSourceName) // Change to "mysql" for MySQL
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(idleTimeout)

	return &DBPool{
		pool:       db,
		maxOpen:    maxOpen,
		maxIdle:    maxIdle,
		idleTimeout: idleTimeout,
	}, nil
}

func (p *DBPool) GetConnection() (*sql.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.pool.Conn(context.Background())
}

func (p *DBPool) Close() error {
	return p.pool.Close()
}
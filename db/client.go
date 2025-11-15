package db

import (
	"context"
	"database/sql"
)

// DBClient defines the interface for unified database operations.
type DBClient interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Find(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	Insert(ctx context.Context, query string, args ...interface{}) (int64, error)
	Update(ctx context.Context, query string, args ...interface{}) (int64, error)
	Close() error
}
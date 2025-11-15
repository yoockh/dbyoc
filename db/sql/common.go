package sql

import (
	"database/sql"
	"fmt"
)

// Common types and functions for SQL database implementations

// DBConfig holds the configuration for a database connection.
type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// OpenDB opens a new database connection based on the provided configuration.
func OpenDB(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) error {
	return db.Close()
}
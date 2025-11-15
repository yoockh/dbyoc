package sql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	*sql.DB
}

func NewMySQLDB(dataSourceName string) (*MySQLDB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &MySQLDB{db}, nil
}

func (m *MySQLDB) Close() error {
	return m.DB.Close()
}

// Additional MySQL specific methods can be added here, such as transaction handling, etc.
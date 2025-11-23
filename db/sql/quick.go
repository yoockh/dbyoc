package sql

import (
	"github.com/yoockh/dbyoc/config"
)

// QuickPostgres creates Postgres client from DATABASE_URL env only
func QuickPostgres() (*PostgresClient, error) {
	cfg, err := config.QuickPostgresConfig()
	if err != nil {
		return nil, err
	}

	return NewPostgresClient(*cfg)
}

// QuickMySQL creates MySQL client from DATABASE_URL env only
func QuickMySQL() (*MySQLDB, error) {
	cfg, err := config.QuickPostgresConfig() // reuse for mysql
	if err != nil {
		return nil, err
	}

	return NewMySQLDB(cfg.URL)
}

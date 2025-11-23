package nosql

import (
	"fmt"

	"github.com/yoockh/dbyoc/config"
)

// Common types and functions for NoSQL database implementations

type NoSQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func NewNoSQLConfig(host string, port int, username, password, database string) *NoSQLConfig {
	return &NoSQLConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}

// For example, a function to validate the NoSQL configuration
func (config *NoSQLConfig) Validate() error {
	if config.Host == "" || config.Port == 0 || config.Database == "" {
		return fmt.Errorf("invalid NoSQL configuration: %+v", config)
	}
	return nil
}

// QuickMongo creates MongoDB client from MONGO_URI env only
func QuickMongo(collection string) (*MongoDBClient, error) {
	cfg, err := config.QuickMongoConfig()
	if err != nil {
		return nil, err
	}

	return NewMongoDBClient(cfg.URI, cfg.Database, collection)
}

package nosql

import (
	"fmt"
	"os"

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

// QuickMongo loads MONGO_URI from env and returns config
// Usage:
//
//	cfg, err := nosql.QuickMongo()
//	if err != nil { log.Fatal(err) }
//	client, err := nosql.NewMongoDBClient(cfg.MongoDB.URI, "mydb", "mycollection")
func QuickMongo() (*config.Config, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = os.Getenv("MONGODB_URI")
	}

	if uri == "" {
		return nil, fmt.Errorf("MONGO_URI or MONGODB_URI environment variable is required")
	}

	return &config.Config{
		MongoDB: config.MongoConfig{
			URI: uri,
		},
		Logger: config.LoggerConfig{
			Level: "info",
		},
	}, nil
}

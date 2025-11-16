# DBYOC - Database Bring Your Own Connection

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/badge/release-v1.2.0-blue.svg)]()

![dbyoc logo](assets/dbyoc.png)

DBYOC is a Go module that simplifies database connections by providing a unified interface for both SQL and NoSQL databases. This module comes with flexible configuration, automatic retry mechanisms, connection pooling, migration support, logging, and metrics tracking.

## Features

- **Unified Interface**: Consistent interface for various database types
- **Flexible Configuration**: Load configuration from environment variables, JSON, or YAML
- **Auto Retry & Reconnect**: Handle transient errors with automatic retry logic
- **Connection Pooling**: Efficient database connection management
- **Built-in Migration**: Manage database schema changes easily
- **Integrated Logging**: Automatic logging for queries and errors using Logrus
- **Metrics Tracking**: Monitor connection and query performance

## Installation

```bash
go get github.com/yoockh/dbyoc
```

## Quick Start

### 1. Create Configuration File

Create a `config.yaml` file:

```yaml
database:
  url: "postgres://user:pass@localhost:5432/mydb?sslmode=disable"  # or use individual fields
  # host: localhost
  # port: 5432
  # user: postgres
  # password: secret
  # database: mydb
  # sslmode: disable
  max_retries: 3
  max_pool_size: 10

mongodb:
  uri: mongodb://localhost:27017
  database: mydb
  timeout: 30

redis:
  url: "redis://:password@localhost:6379/0"  # or use individual fields
  # addr: localhost:6379
  # password: ""
  # db: 0

logger:
  level: info
```

### 2. PostgreSQL Connection

```go
package main

import (
    "log"
    
    "github.com/yoockh/dbyoc/config"
    "github.com/yoockh/dbyoc/db/sql"
    "github.com/yoockh/dbyoc/logger"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Initialize logger
    logger.Init(cfg.Logger.Level)
    
    // Connect to PostgreSQL
    pgClient, err := sql.NewPostgresClient(cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer pgClient.Close()
    
    // Execute query
    rows, err := pgClient.Query("SELECT id, name FROM users LIMIT 10")
    if err != nil {
        log.Fatal("Query failed:", err)
    }
    defer rows.Close()
    
    // Process results
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            log.Fatal("Scan failed:", err)
        }
        logger.GetLogger().Infof("User: ID=%d, Name=%s", id, name)
    }
}
```

### 3. MongoDB Connection

```go
package main

import (
    "context"
    "log"
    
    "github.com/yoockh/dbyoc/config"
    "github.com/yoockh/dbyoc/db/nosql"
    "github.com/yoockh/dbyoc/logger"
    "go.mongodb.org/mongo-driver/bson"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Initialize logger
    logger.Init(cfg.Logger.Level)
    
    // Connect to MongoDB
    mongoClient, err := nosql.NewMongoDBClient(
        cfg.MongoDB.URI,
        cfg.MongoDB.Database,
        "users",
    )
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer mongoClient.Close()
    
    // Insert document
    user := bson.M{"name": "John Doe", "email": "john@example.com"}
    if err := mongoClient.Insert(user); err != nil {
        log.Fatal("Insert failed:", err)
    }
    
    // Find documents
    cursor, err := mongoClient.Find(bson.M{})
    if err != nil {
        log.Fatal("Find failed:", err)
    }
    defer cursor.Close(context.Background())
    
    for cursor.Next(context.Background()) {
        var result bson.M
        if err := cursor.Decode(&result); err != nil {
            log.Fatal("Decode failed:", err)
        }
        logger.GetLogger().Infof("Document: %+v", result)
    }
}
```

### 4. Redis Connection

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/yoockh/dbyoc/config"
    "github.com/yoockh/dbyoc/db/nosql"
    "github.com/yoockh/dbyoc/logger"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Initialize logger
    logger.Init(cfg.Logger.Level)
    log := logger.GetLogger()
    
    // Connect to Redis
    redisClient := nosql.NewRedisClient(
        cfg.Redis.Addr,
        cfg.Redis.Password,
        cfg.Redis.DB,
        log,
    )
    defer redisClient.Close()
    
    ctx := context.Background()
    
    // Ping connection
    if err := redisClient.Ping(ctx); err != nil {
        log.Fatal("Ping failed:", err)
    }
    
    // Set value with expiration
    if err := redisClient.Set(ctx, "user:1", "John Doe", 5*time.Minute); err != nil {
        log.Fatal("Set failed:", err)
    }
    
    // Get value
    val, err := redisClient.Get(ctx, "user:1")
    if err != nil {
        log.Fatal("Get failed:", err)
    }
    
    logger.GetLogger().Infof("Value: %s", val)
}
```

### 5. Database Migration

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    "github.com/yoockh/dbyoc/config"
    "github.com/yoockh/dbyoc/migration"
    _ "github.com/lib/pq"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // Open database connection
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host, cfg.Database.Port, cfg.Database.User, 
        cfg.Database.Password, cfg.Database.Database, cfg.Database.SSLMode)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Run migrations
    migrator := migration.NewMigrator(db)
    if err := migrator.Migrate("./migration/files"); err != nil {
        log.Fatal("Migration failed:", err)
    }
    
    log.Println("Migration completed successfully")
}
```

Create migration file `migration/files/001_create_users_table.sql`:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 6. Metrics Monitoring

```go
package main

import (
    "log"
    "time"
    
    "github.com/yoockh/dbyoc/metrics"
)

func main() {
    // Create metrics collector
    collector := metrics.NewMetricsCollector()
    
    // Track connections and queries
    collector.IncrementConnectionCount()
    collector.IncrementQueryCount()
    
    time.Sleep(100 * time.Millisecond)
    
    // Retrieve metrics
    queryCount, connCount, uptime := collector.GetMetrics()
    log.Printf("Queries: %d, Connections: %d, Uptime: %v", 
        queryCount, connCount, uptime)
}
```

## Project Structure

```
dbyoc/
├── config/          # Configuration management (Viper)
├── db/             
│   ├── sql/        # PostgreSQL, MySQL implementations
│   └── nosql/      # MongoDB, Redis implementations
├── migration/      # Database migration tools
├── logger/         # Logging (Logrus)
├── metrics/        # Metrics tracking
└── utils/          # Retry logic, helpers
```

## Configuration

DBYOC uses **Viper** for configuration management. You can load configuration from:

- YAML file
- JSON file
- Environment variables (standard naming)

### Using Config File (config.yaml)

```yaml
database:
  url: "postgres://user:pass@localhost:5432/mydb?sslmode=disable"  # or use individual fields
  # host: localhost
  # port: 5432
  # user: postgres
  # password: secret
  # database: mydb
  # sslmode: disable
  max_retries: 3
  max_pool_size: 10

mongodb:
  uri: mongodb://localhost:27017
  database: mydb
  timeout: 30

redis:
  url: "redis://:password@localhost:6379/0"  # or use individual fields
  # addr: localhost:6379
  # password: ""
  # db: 0

logger:
  level: info
```

### Using Environment Variables

Simple and standard naming:

```bash
# PostgreSQL - using URL
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"

# Or individual fields
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=secret
export DATABASE_NAME=mydb
export DATABASE_SSLMODE=disable

# MongoDB
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE=mydb

# Redis - using URL
export REDIS_URL="redis://:password@localhost:6379/0"

# Or individual fields
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=secret
export REDIS_DB=0

# Logger
export LOG_LEVEL=debug
```

### Load Config in Code

```go
// Load from config file (config.yaml)
cfg, err := config.LoadConfig()

// Or load purely from environment variables
cfg, err := config.LoadFromEnv()
```

## External Dependencies

DBYOC leverages proven, production-ready external libraries:

| Library | Purpose | Repository |
|---------|---------|------------|
| **Logrus** | Structured logging | [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) |
| **Viper** | Configuration management | [github.com/spf13/viper](https://github.com/spf13/viper) |
| **Mapstructure** | Map to struct conversion | [github.com/mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) |
| **PostgreSQL Driver** | PostgreSQL connectivity | [github.com/lib/pq](https://github.com/lib/pq) |
| **MySQL Driver** | MySQL connectivity | [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) |
| **MongoDB Driver** | MongoDB connectivity | [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver) |
| **Redis Client** | Redis connectivity | [github.com/go-redis/redis](https://github.com/go-redis/redis) |

## Contributing

Contributions are welcome. To contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Guidelines

- Follow Go conventions and best practices
- Add tests for new features
- Update documentation as needed
- Use Conventional Commits for commit messages

## Changelog
### v1.3.0 (Latest)
- Server wrapper now supports using Echo directly as the HTTP handler (Echo implements `http.Handler`).
- Introduced `StartWithSignals()` to:
  - Run the HTTP or HTTPS server.
  - Handle OS signals (SIGINT, SIGTERM).
  - Perform graceful shutdown using `cfg.Server.ShutdownTimeout`.
- Example usage added showing how to integrate Echo and the server wrapper in `main.go`.

### v1.2.5
[Implement HTTP server with graceful shutdown](https://github.com/yoockh/dbyoc/commit/a9ff77063da4b4115a61648100ffded4b35a8045)

### v1.2.0
- Config
  - LoadFromEnv now ensures environment variables are read when it runs standalone (viper.AutomaticEnv added).
  - Added Config.Validate() to enforce minimal required settings:
    - `database.type` is required.
    - Either `database.url`, or `database.host` + `database.port` + `database.database` is required.
  - Kept existing environment variable naming (DATABASE_*, MONGODB_*, REDIS_*, LOG_LEVEL) for backward compatibility.
  - Improved README/docs to clarify env names, YAML lookup paths, and recommended usage.
- Documentation
  - Root README and per-package READMEs updated for consistent style and full examples.
  - Added a detailed Changelog section.
- Utilities & Clients
  - Clarified how to compose Retry and Backoff (examples in utils README).
  - Ensured Redis client exposes Ping, Reconnect, Set/Get, RetryOperation behavior clearly.
  - MongoDB client performs a connection check (Ping) on initialization and provides simple CRUD wrappers.
- Stability & Tests
  - Minor bug fixes and code tidying across packages.
  - Added recommendations for unit tests around config loading and validation.
- Upgrade notes:
  - Call `cfg.Validate()` after loading configuration and handle validation errors before initializing clients.
  - No breaking changes to public API or environment variable names.

  ### v1.1.0
  [Add FindOne and Delete methods to MongoDBClient](https://github.com/yoockh/dbyoc/commit/afa3c9e358d8714139c1cec4d03b21469e0061c5)

  ### v1.0.0
- **Unified Interface**: Consistent interface for various database types
- **Flexible Configuration**: Load configuration from environment variables, JSON, or YAML
- **Auto Retry & Reconnect**: Handle transient errors with automatic retry logic
- **Connection Pooling**: Efficient database connection management
- **Built-in Migration**: Manage database schema changes easily
- **Integrated Logging**: Automatic logging for queries and errors using Logrus
- **Metrics Tracking**: Monitor connection and query performance

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Author

**Aisiya Qutwatunnada**

---

**Note**: DBYOC is designed to simplify database connections, not complicate them. If you encounter any issues or have suggestions, please create an issue or pull request.

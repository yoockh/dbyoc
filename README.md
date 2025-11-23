# DBYOC - Database Bring Your Own Connection

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/badge/release-v1.3.6-blue.svg)]()

![dbyoc logo](assets/dbyoc.png)

DBYOC is a Go module that simplifies database connections by providing a unified interface for both SQL and NoSQL databases. This module comes with flexible configuration, automatic retry mechanisms, connection pooling, migration support, logging, and metrics tracking.

## Features

- **One-Line Setup**: Connect with just environment variable
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

## Super Quick Start (Recommended)

### MongoDB - One Line Setup

```bash
export MONGO_URI="mongodb://localhost:27017/mydb"
```

```go
package main

import (
    "context"
    "log"
    
    "github.com/yoockh/dbyoc/db/nosql"
    "go.mongodb.org/mongo-driver/bson"
)

func main() {
    // ONE LINE!
    client, err := nosql.QuickMongo("users")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Insert
    client.Insert(bson.M{"name": "John", "email": "john@example.com"})
    
    // Find
    cursor, _ := client.Find(bson.M{})
    defer cursor.Close(context.Background())
    
    for cursor.Next(context.Background()) {
        var result bson.M
        cursor.Decode(&result)
        log.Printf("User: %+v", result)
    }
}
```

### PostgreSQL - One Line Setup

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
```

```go
package main

import (
    "log"
    "github.com/yoockh/dbyoc/db/sql"
)

func main() {
    // ONE LINE!
    db, err := sql.QuickPostgres()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    rows, _ := db.Query("SELECT id, name FROM users")
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        log.Printf("User: %d - %s", id, name)
    }
}
```

### Redis - One Line Setup

```bash
export REDIS_URL="redis://:password@localhost:6379/0"
```

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/yoockh/dbyoc/db/nosql"
)

func main() {
    // ONE LINE!
    redis, err := nosql.QuickRedis()
    if err != nil {
        log.Fatal(err)
    }
    defer redis.Close()
    
    ctx := context.Background()
    redis.Set(ctx, "key", "value", 5*time.Minute)
    
    val, _ := redis.Get(ctx, "key")
    log.Printf("Value: %s", val)
}
```

## Environment Variables

### Minimal (Recommended)

```bash
# MongoDB
export MONGO_URI="mongodb://localhost:27017/mydb"

# PostgreSQL
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb"

# Redis
export REDIS_URL="redis://:password@localhost:6379/0"
```

### Optional Settings

```bash
# MongoDB
export MONGO_DATABASE="mydb"
export MONGO_TIMEOUT=30

# PostgreSQL
export DATABASE_MAX_RETRIES=3
export DATABASE_MAX_POOL_SIZE=10

# Redis
export REDIS_PASSWORD="secret"
export REDIS_DB=0
```

## Advanced Usage (Optional)

If you need more control, use config file:

```yaml
database:
  url: "postgres://user:pass@localhost:5432/mydb"
  max_pool_size: 10

mongodb:
  uri: mongodb://localhost:27017
  database: mydb

redis:
  url: "redis://:password@localhost:6379/0"
```

```go
cfg, err := config.LoadConfig()
```

## Project Structure

```
dbyoc/
├── config/          # Configuration management
├── db/             
│   ├── sql/        # PostgreSQL, MySQL
│   └── nosql/      # MongoDB, Redis
├── migration/      # Database migrations
├── logger/         # Logging
├── metrics/        # Metrics tracking
└── utils/          # Utilities
```

## External Dependencies

| Library | Purpose | Repository |
|---------|---------|------------|
| **Logrus** | Structured logging | [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) |
| **Viper** | Configuration management | [github.com/spf13/viper](https://github.com/spf13/viper) |
| **Mapstructure** | Map to struct conversion | [github.com/mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) |
| **PostgreSQL Driver** | PostgreSQL connectivity | [github.com/lib/pq](https://github.com/lib/pq) |
| **MySQL Driver** | MySQL connectivity | [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) |
| **MongoDB Driver** | MongoDB connectivity | [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver) |
| **Redis Client** | Redis connectivity | [github.com/go-redis/redis](https://github.com/go-redis/redis) |

## Changelog

### v1.3.7 (Latest)
- Added Quick Start helpers: `QuickMongo()`, `QuickPostgres()`, `QuickRedis()`
- MongoDB works with just `MONGO_URI` environment variable
- PostgreSQL works with just `DATABASE_URL` environment variable
- Redis works with just `REDIS_URL` environment variable
- No complex configuration needed for basic usage
- Backward compatible with existing configuration methods

### v1.3.6
- Fixed module import resolution issues
- Improved graceful shutdown flow
- Corrected timeout handling

## Contributing

Contributions are welcome. To contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Author

**Aisiya Qutwatunnada**

---

**Note**: DBYOC is designed to simplify database connections. Just set one environment variable and you're ready to go.

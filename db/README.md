# db — Database helpers for DBYOC

This folder contains database helpers, clients, and utilities used by DBYOC. The goal is to provide a unified, lightweight abstraction over SQL and NoSQL databases with practical helpers for connection pooling, retry logic, and simple client wrappers for common databases (Postgres, MySQL, MongoDB, Redis).

- Language: Go
- Location: `./db`

## Key features
- A unified DBClient interface for common database operations.
- SQL helpers and clients:
  - Postgres client with connection pool configuration and retry helpers.
  - MySQL client wrapper with sensible connection pooling defaults.
  - Generic utilities to open and close SQL connections.
- NoSQL clients:
  - MongoDB client wrapping the official Mongo driver.
  - Redis client (go-redis) with logging, ping/reconnect helpers and retry support.
- Database connection pool helper (DBPool) for managing long-lived connections.

## Package layout
- `db/`
  - `client.go` — DBClient interface defining common operations.
  - `pool.go` — DBPool: a thin wrapper for sql.DB connection pooling.
- `db/sql/`
  - `common.go` — SQL helpers (DBConfig, OpenDB, CloseDB).
  - `postgres.go` — PostgresClient with pool and retry helpers.
  - `mysql.go` — MySQLDB wrapper with pool settings.
- `db/nosql/`
  - `common.go` — NoSQL configuration type and helpers.
  - `mongo.go` — MongoDBClient wrapper.
  - `redis.go` — RedisClient wrapper (go-redis) with logging and reconnect.

## Interface (db/client.go)
DBClient defines a minimal unified interface:
- Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
- Find(ctx context.Context, query string, args ...interface{}) (interface{}, error)
- Insert(ctx context.Context, query string, args ...interface{}) (int64, error)
- Update(ctx context.Context, query string, args ...interface{}) (int64, error)
- Close() error

This interface is intended for higher-level code that needs to operate against different database backends with a consistent API.

## Examples

Postgres example (uses config.DatabaseConfig from the repo):
```go
package main

import (
	"log"

	"github.com/yoockh/dbyoc/config"
	sqlpkg "github.com/yoockh/dbyoc/db/sql"
)

func main() {
	cfg := config.DatabaseConfig{
		Host:        "localhost",
		Port:        5432,
		User:        "postgres",
		Password:    "secret",
		Database:    "mydb",
		SSLMode:     "disable",
		MaxPoolSize: 20,
		MaxRetries:  3,
	}

	client, err := sqlpkg.NewPostgresClient(cfg)
	if err != nil {
		log.Fatalf("failed to create postgres client: %v", err)
	}
	defer client.Close()

	rows, err := client.Query("SELECT id, name FROM users WHERE active = $1", true)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}
	defer rows.Close()
	// iterate rows...
}
```

MySQL example:
```go
db, err := sqlpkg.NewMySQLDB("user:pass@tcp(localhost:3306)/dbname")
if err != nil { /* handle error */ }
defer db.Close()
// use db (MySQLDB embeds *sql.DB)
```

MongoDB example:
```go
package main

import (
	"log"

	"github.com/yoockh/dbyoc/db/nosql"
)

func main() {
	mClient, err := nosql.NewMongoDBClient("mongodb://user:pass@localhost:27017", "mydb", "collection")
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer mClient.Close()

	if err := mClient.Insert(map[string]interface{}{"name": "alice"}); err != nil {
		log.Fatalf("insert failed: %v", err)
	}
}
```

Redis example:
```go
package main

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/yoockh/dbyoc/config"
	"github.com/yoockh/dbyoc/db/nosql"
)

func main() {
	logger := logrus.New()
	cfg := config.RedisConfig{
		Addr: "localhost:6379",
		DB:   0,
	}

	r := nosql.NewRedisClient(cfg, logger)
	defer r.Close()

	ctx := context.Background()
	if err := r.Set(ctx, "foo", "bar", 0); err != nil {
		log.Fatalf("redis set: %v", err)
	}

	val, _ := r.Get(ctx, "foo")
	log.Println("value:", val)
}
```

DBPool example:
```go
package main

import (
	"time"

	"github.com/yoockh/dbyoc/db"
)

func main() {
	pool, err := db.NewDBPool("host=localhost port=5432 user=... dbname=...", 50, 25, 30*time.Minute)
	if err != nil { /* handle */ }
	defer pool.Close()

	conn, err := pool.GetConnection()
	if err != nil { /* handle */ }
	defer conn.Close()

	// Use conn for connection-scoped operations
	_ = conn
}
```

## Notes and recommendations
- The SQL helpers assume usage of standard Go drivers (e.g., lib/pq for Postgres, go-sql-driver/mysql for MySQL). Make sure the appropriate driver is imported in your application.
- Redis client depends on a repository-level `config.RedisConfig` type and a logger (logrus). Adapt as needed for your environment.
- Mongo client uses the official mongo-driver and provides basic CRUD wrappers; expand it as required for transactions and advanced options.
- The Postgres client includes a RetryQuery helper — tune retry settings via the repository config struct.

## Contributing
Contributions, bug reports, and improvements are welcome. Follow the project contribution guidelines in the repository root.

## License
See the repository LICENSE file for details.

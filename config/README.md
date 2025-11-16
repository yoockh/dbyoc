# config — Configuration loader for DBYOC

This package provides a small, practical configuration loader for DBYOC. It supports loading configuration from YAML files and environment variables (via Viper) and decodes into typed Go structs using mapstructure.

- Language: Go
- Location: `./config`

## Purpose

Centralize application configuration for database and logging settings. The package is intentionally simple and aims to:

- Load configuration from a YAML file (with optional custom path)
- Allow environment variables to override file values
- Provide a lightweight Validate method to ensure minimal required configuration

## Files

- `config.go` — main types and loader functions:
  - Config, DatabaseConfig, MongoConfig, RedisConfig, LoggerConfig
  - LoadConfig(configPath ...string) (*Config, error)
  - LoadFromEnv() (*Config, error)
  - (Config).Validate() error

## Types

- Config
  - Database DatabaseConfig `mapstructure:"database"`
  - MongoDB  MongoConfig    `mapstructure:"mongodb"`
  - Redis    RedisConfig    `mapstructure:"redis"`
  - Logger   LoggerConfig   `mapstructure:"logger"`

- DatabaseConfig
  - Type, Host, Port, User, Password, Database, SSLMode
  - MaxRetries, MaxPoolSize
  - URL (full connection string support)

- MongoConfig
  - URI, Database, Timeout (int; currently treated as seconds)

- RedisConfig
  - Addr, Password, DB, URL

- LoggerConfig
  - Level

## Key functions

- LoadConfig(configPath ...string) (*Config, error)  
  Loads configuration from a YAML file (default lookup: `./config/config.yaml` or `./config.yaml`) and enables environment variable overrides via Viper's AutomaticEnv. If a configPath is provided it will be used as the file path.

- LoadFromEnv() (*Config, error)  
  Builds the Config struct purely from environment variables. This function calls viper.AutomaticEnv() so it works even if LoadConfig() has not been invoked.

- (cfg *Config) Validate() error  
  Performs conservative validation:
  - `database.type` is required.
  - Either `database.url` must be provided, or `database.host`, `database.port`, and `database.database` must be present.
  - If `mongodb.uri` is set, `mongodb.database` is recommended (treated as error to encourage explicitness).
  - Redis checks are lenient by default (do not fail if Redis is unused).

Call Validate() after loading config to ensure required settings are present.

## Environment variables

The package reads environment variables using the names below (upper-case). These variables are read directly by Viper using GetString/GetInt in LoadFromEnv:

- Database:
  - DATABASE_TYPE
  - DATABASE_HOST
  - DATABASE_PORT
  - DATABASE_USER
  - DATABASE_PASSWORD
  - DATABASE_NAME
  - DATABASE_SSLMODE
  - DATABASE_URL
  - DATABASE_MAX_RETRIES
  - DATABASE_MAX_POOL_SIZE

- MongoDB:
  - MONGODB_URI
  - MONGODB_DATABASE
  - MONGODB_TIMEOUT

- Redis:
  - REDIS_ADDR
  - REDIS_PASSWORD
  - REDIS_DB
  - REDIS_URL

- Logger:
  - LOG_LEVEL

Note: LoadConfig enables viper.AutomaticEnv(), so those same environment variables will override YAML file values when present.

## Examples

YAML example (example `config/config.yaml`):
```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "secret"
  database: "mydb"
  sslmode: "disable"
  max_retries: 3
  max_pool_size: 20
  url: ""

mongodb:
  uri: "mongodb://user:pass@localhost:27017"
  database: "mydb"
  timeout: 10

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  url: ""

logger:
  level: "info"
```

Load from file (with optional path):
```go
cfg, err := config.LoadConfig() // searches ./ and ./config
if err != nil { /* handle */ }
if err := cfg.Validate(); err != nil { /* handle invalid config */ }
```

Load purely from environment:
```go
cfg, err := config.LoadFromEnv()
if err != nil { /* handle */ }
if err := cfg.Validate(); err != nil { /* handle invalid config */ }
```

## Notes and recommendations

- LoadFromEnv now calls viper.AutomaticEnv() so it reliably reads environment variables even if LoadConfig was never called.
- Validation is intentionally conservative. Tweak Validate() to be stricter (e.g., require Redis) if your application depends on those services.
- Consider adding defaults for common DB ports or log level if desired.
- Consider converting MongoConfig.Timeout to type time.Duration and using mapstructure StringToTimeDurationHookFunc for parsing duration strings (e.g., "10s", "1m").
- If you want environment variables to automatically map to nested keys (e.g., DATABASE_HOST -> database.host for Unmarshal), consider using a local viper instance and SetEnvKeyReplacer(strings.NewReplacer(".", "_")) — currently the code reads env vars by explicit uppercase keys in LoadFromEnv for clarity and backward compatibility.

## Testing suggestions

- Unit tests for LoadFromEnv: set environment variables in test and assert config values.
- Unit tests for LoadConfig: provide temporary YAML files and assert fields; test env override behavior.
- Tests for Validate(): check expected failures for missing required fields.

## Contribution and License

Contributions and improvements are welcome. Follow repository contribution guidelines. See the repository LICENSE file for license details.

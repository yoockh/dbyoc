# Config Package

The `config` package provides a clean way to manage your application's configuration in Go.  
It supports loading configuration from **YAML files** or **environment variables**.

---

## Package Overview

The package includes the following structs:

- **Config**: Main struct holding all configurations (Database, MongoDB, Redis, Logger)
- **DatabaseConfig**: Database connection settings (supports URL or individual fields)
- **MongoConfig**: MongoDB connection settings
- **RedisConfig**: Redis connection settings (supports URL or individual fields)
- **LoggerConfig**: Logger settings (e.g., log level)

### Key Functions

1. **LoadConfig(configPath ...string) (*Config, error)**  
   Loads configuration from a YAML file.  
   - If a path is provided, loads from that file  
   - If no path is provided, searches `./config/config.yaml` or `./config.yaml`  
   - Automatically binds to environment variables if present  
   - Returns a fully populated `Config` struct  

2. **LoadFromEnv() (*Config, error)**  
   Loads configuration purely from environment variables.  

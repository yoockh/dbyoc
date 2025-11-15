# DBYOC
## Database Bring Your Own Connection

DBYOC is a Go module designed to simplify database connections by providing a unified interface for both SQL and NoSQL databases. It offers flexible configuration options, automatic retry and reconnect mechanisms, connection pooling, built-in migration support, integrated logging, and metrics tracking.

### Features
- **Unified Database Connection**: Interact with various databases using a consistent interface.
- **Flexible Configuration**: Load configurations from environment variables, JSON, or YAML files.
- **Automatic Retry and Reconnect**: Handle transient errors gracefully with built-in retry logic.
- **Connection Pooling**: Efficiently manage database connections to optimize resource usage.
- **Built-in Migration**: Easily manage database schema changes with migration support.
- **Integrated Logging**: Automatically log queries and errors for better observability.
- **Metrics Tracking**: Monitor connection and query performance with built-in metrics.

### Project Structure
```
dbyoc/
├── config/                  # Configuration management
│   ├── config.go            # Configuration structure and loading methods
│   ├── loader.go            # Logic for loading configuration files
│   └── types.go             # Types and constants for configuration
│
├── db/                      # Database connection and handling
│   ├── client.go            # Unified DBClient interface
│   ├── pool.go              # Connection pooling management
│   ├── sql/                 # SQL specific implementations
│   │   ├── postgres.go      # PostgreSQL connection logic
│   │   ├── mysql.go         # MySQL connection logic
│   │   └── common.go        # Common SQL functions and types
│   └── nosql/               # NoSQL specific implementations
│       ├── mongo.go         # MongoDB connection logic
│       ├── redis.go         # Redis connection logic
│       └── common.go        # Common NoSQL functions and types
│
├── migration/               # Database migration management
│   ├── migration.go         # Main migration functionality
│   ├── migrator.go          # Logic for executing migration scripts
│   └── files/               # Directory for migration files
│       └── .gitkeep         # Keep migration files directory in version control
│
├── logger/                  # Logging helper
│   ├── logger.go            # Integrated logging functionality
│   └── interface.go         # Logging interface and methods
│
├── metrics/                 # Metrics tracking
│   ├── metrics.go           # Metrics functionality
│   ├── collector.go         # Metrics data collection and processing
│   └── types.go             # Types and constants for metrics
│
├── utils/                   # Utility functions
│   ├── retry.go             # Retry logic and backoff strategies
│   ├── backoff.go           # Exponential backoff functions
│   └── helpers.go           # Various utility functions
│
├── example/                 # Example usage of the package
│   ├── main.go              # General example implementation
│   ├── postgres_example.go   # Example with PostgreSQL
│   └── mongo_example.go      # Example with MongoDB
│
├── go.mod                   # Module definition and dependencies
├── go.sum                   # Dependency checksums
├── README.md                # Project documentation
└── LICENSE                  # Licensing information
```

### Installation
To install the DBYOC module, use the following command:

```
go get github.com/yourusername/dbyoc
```

### Usage
Refer to the examples in the `example` directory for practical implementations of the DBYOC package with different databases. 

### License
This project is licensed under the MIT License. See the LICENSE file for more details.
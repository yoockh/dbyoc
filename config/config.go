package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	MongoDB  MongoConfig    `mapstructure:"mongodb"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Type        string `mapstructure:"type"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	SSLMode     string `mapstructure:"sslmode"`
	MaxRetries  int    `mapstructure:"max_retries"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
	// Support connection string/URL
	URL string `mapstructure:"url"`
}

type MongoConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	// Support Redis URL
	URL string `mapstructure:"url"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// ServerConfig holds HTTP server related settings.
type ServerConfig struct {
	// Either set Addr directly (host:port) or set Host + Port
	Addr string `mapstructure:"addr"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	// timeouts in seconds (simple int to keep config YAML/env friendly)
	ReadTimeout     int  `mapstructure:"read_timeout"`
	WriteTimeout    int  `mapstructure:"write_timeout"`
	ShutdownTimeout int  `mapstructure:"shutdown_timeout"`
	TLS             bool `mapstructure:"tls"`
	CertFile        string `mapstructure:"cert_file"`
	KeyFile         string `mapstructure:"key_file"`
}

// Address returns the effective listen address for the server.
func (s *ServerConfig) Address() string {
	if s == nil {
		return ":8080"
	}
	if s.Addr != "" {
		return s.Addr
	}
	host := s.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := s.Port
	if port == 0 {
		port = 8080
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func LoadConfig(configPath ...string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if len(configPath) > 0 {
		viper.SetConfigFile(configPath[0])
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	// Auto bind environment variables (so env vars can override file keys)
	viper.AutomaticEnv()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		// If no config file, that's okay — we'll rely on env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config, viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc())); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// LoadFromEnv loads config purely from environment variables.
//
// NOTE: we call viper.AutomaticEnv() here to ensure viper reads env vars even if
// LoadConfig wasn't called beforehand.
func LoadFromEnv() (*Config, error) {
	viper.AutomaticEnv()

	return &Config{
		Database: DatabaseConfig{
			Type:        viper.GetString("DATABASE_TYPE"),
			Host:        viper.GetString("DATABASE_HOST"),
			Port:        viper.GetInt("DATABASE_PORT"),
			User:        viper.GetString("DATABASE_USER"),
			Password:    viper.GetString("DATABASE_PASSWORD"),
			Database:    viper.GetString("DATABASE_NAME"),
			SSLMode:     viper.GetString("DATABASE_SSLMODE"),
			URL:         viper.GetString("DATABASE_URL"),
			MaxRetries:  viper.GetInt("DATABASE_MAX_RETRIES"),
			MaxPoolSize: viper.GetInt("DATABASE_MAX_POOL_SIZE"),
		},
		MongoDB: MongoConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
			Timeout:  viper.GetInt("MONGODB_TIMEOUT"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			URL:      viper.GetString("REDIS_URL"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		Server: ServerConfig{
			Addr:            viper.GetString("SERVER_ADDR"),
			Host:            viper.GetString("SERVER_HOST"),
			Port:            viper.GetInt("SERVER_PORT"),
			ReadTimeout:     viper.GetInt("SERVER_READ_TIMEOUT"),
			WriteTimeout:    viper.GetInt("SERVER_WRITE_TIMEOUT"),
			ShutdownTimeout: viper.GetInt("SERVER_SHUTDOWN_TIMEOUT"),
			TLS:             viper.GetBool("SERVER_TLS"),
			CertFile:        viper.GetString("SERVER_CERT_FILE"),
			KeyFile:         viper.GetString("SERVER_KEY_FILE"),
		},
	}, nil
}

// Validate checks required configuration values and returns an error describing
// any missing or inconsistent settings. This is intentionally conservative:
// it enforces minimal requirements for the Database config (which most apps need),
// and performs light checks for other subsystems.
func (c *Config) Validate() error {
	// Database.Type is important for choosing driver; require it.
	if c.Database.Type == "" {
		return fmt.Errorf("database.type is required")
	}

	// Database must have either a full URL or host+port+database set.
	if c.Database.URL == "" {
		if c.Database.Host == "" {
			return fmt.Errorf("database.host is required when database.url is not provided")
		}
		if c.Database.Port == 0 {
			return fmt.Errorf("database.port is required when database.url is not provided")
		}
		if c.Database.Database == "" {
			// database name is required for most DB drivers
			return fmt.Errorf("database.database (name) is required when database.url is not provided")
		}
	}

	// MongoDB: if URI provided, database is recommended (but not strictly required).
	if c.MongoDB.URI != "" && c.MongoDB.Database == "" {
		// warn as error to force explicitness; change to non-fatal if you prefer
		return fmt.Errorf("mongodb.database is recommended when MONGODB_URI is set")
	}

	// Redis: require at least one connection method
	if c.Redis.Addr == "" && c.Redis.URL == "" {
		// not fatal if your app doesn't use Redis; caller can decide whether Redis is required.
		// We return an error only if caller expects Redis — but since this function can't know that,
		// we keep it lenient and do not return an error here. Uncomment the next line to enforce.
		// return fmt.Errorf("redis.addr or redis.url must be provided")
	}

	// Server: if TLS enabled must provide cert and key
	if c.Server.TLS {
		if c.Server.CertFile == "" || c.Server.KeyFile == "" {
			return fmt.Errorf("server.tls is enabled but cert_file or key_file is missing")
		}
	}

	// Logger: set a sensible default externally; we don't treat empty log level as an error.
	return nil
}

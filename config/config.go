package config

import (
    "fmt"
    "log"
    "strings"

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
    URL         string `mapstructure:"url"`
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
    URL      string `mapstructure:"url"`
}

type LoggerConfig struct {
    Level string `mapstructure:"level"`
}

type ServerConfig struct {
    Addr            string `mapstructure:"addr"`
    Host            string `mapstructure:"host"`
    Port            int    `mapstructure:"port"`
    ReadTimeout     int    `mapstructure:"read_timeout"`
    WriteTimeout    int    `mapstructure:"write_timeout"`
    ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
    TLS             bool   `mapstructure:"tls"`
    CertFile        string `mapstructure:"cert_file"`
    KeyFile         string `mapstructure:"key_file"`
}

func (s *ServerConfig) Address() string {
    if s == nil {
        return ":8080"
    }
    if s.Addr != "" {
        // normalize leading ":" if only port
        if strings.HasPrefix(s.Addr, ":") && !strings.Contains(s.Addr, "0.0.0.0") && !strings.Contains(s.Addr, "127.0.0.1") {
            return s.Addr
        }
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
    viper.AutomaticEnv()

    _ = viper.ReadInConfig() // ignore not found

    var cfg Config
    if err := viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc())); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    postProcess(&cfg)
    return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
    viper.AutomaticEnv()
    cfg := &Config{
        Database: DatabaseConfig{
            Type:        viper.GetString("DATABASE_TYPE"),
            Host:        viper.GetString("DATABASE_HOST"),
            Port:        viper.GetInt("DATABASE_PORT"),
            User:        viper.GetString("DATABASE_USER"),
            Password:    viper.GetString("DATABASE_PASSWORD"),
            Database:    viper.GetString("DATABASE_NAME"),
            SSLMode:     viper.GetString("DATABASE_SSLMODE"),
            URL:         firstNonEmpty(
                viper.GetString("DATABASE_URL"),
                viper.GetString("MONGO_URI"),
                viper.GetString("MONGODB_URI"),
            ),
            MaxRetries:  viper.GetInt("DATABASE_MAX_RETRIES"),
            MaxPoolSize: viper.GetInt("DATABASE_MAX_POOL_SIZE"),
        },
        MongoDB: MongoConfig{
            URI:      firstNonEmpty(viper.GetString("MONGODB_URI"), viper.GetString("MONGO_URI")),
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
            Level: firstNonEmpty(viper.GetString("LOGGER_LEVEL"), viper.GetString("LOG_LEVEL"), "info"),
        },
        Server: ServerConfig{
            Addr:            firstNonEmpty(viper.GetString("SERVER_ADDRESS"), viper.GetString("SERVER_ADDR")),
            Host:            viper.GetString("SERVER_HOST"),
            Port:            viper.GetInt("SERVER_PORT"),
            ReadTimeout:     viper.GetInt("SERVER_READ_TIMEOUT"),
            WriteTimeout:    viper.GetInt("SERVER_WRITE_TIMEOUT"),
            ShutdownTimeout: viper.GetInt("SERVER_SHUTDOWN_TIMEOUT"),
            TLS:             viper.GetBool("SERVER_TLS"),
            CertFile:        viper.GetString("SERVER_CERT_FILE"),
            KeyFile:         viper.GetString("SERVER_KEY_FILE"),
        },
    }
    postProcess(cfg)
    return cfg, nil
}

func postProcess(c *Config) {
    // If only Mongo is used, set database.type/url for compatibility
    if c.MongoDB.URI != "" && c.Database.URL == "" && c.Database.Host == "" {
        c.Database.URL = c.MongoDB.URI
        if c.Database.Type == "" {
            c.Database.Type = "mongodb"
        }
    }
    // Normalize logger level
    if c.Logger.Level == "" {
        c.Logger.Level = "info"
    }
}

func (c *Config) UsingMongoOnly() bool {
    return c.Database.Type == "mongodb" &&
        c.Database.URL != "" &&
        c.MongoDB.URI != "" &&
        c.Database.Host == "" &&
        c.Database.Port == 0
}

func (c *Config) Validate() error {
    // If pure Mongo usage, minimal checks
    if c.UsingMongoOnly() {
        // Basic sanity
        if !strings.HasPrefix(c.MongoDB.URI, "mongodb") {
            return fmt.Errorf("invalid mongodb uri")
        }
        return nil
    }

    // Generic DB validation
    if c.Database.Type == "" {
        return fmt.Errorf("database.type is required")
    }
    if c.Database.URL == "" {
        if c.Database.Host == "" {
            return fmt.Errorf("database.host is required when database.url is not provided")
        }
        if c.Database.Port == 0 {
            return fmt.Errorf("database.port is required when database.url is not provided")
        }
        if c.Database.Database == "" {
            return fmt.Errorf("database.database (name) is required when database.url is not provided")
        }
    }

    // Mongo optional
    if c.MongoDB.URI != "" && !strings.HasPrefix(c.MongoDB.URI, "mongodb") {
        return fmt.Errorf("mongodb.uri must start with mongodb or mongodb+srv")
    }

    if c.Server.TLS {
        if c.Server.CertFile == "" || c.Server.KeyFile == "" {
            return fmt.Errorf("server.tls enabled but cert/key missing")
        }
    }
    return nil
}

func firstNonEmpty(vals ...string) string {
    for _, v := range vals {
        if strings.TrimSpace(v) != "" {
            return v
        }
    }
    return ""
}

// Simple helper to show final resolved configuration (debug)
func (c *Config) DebugPrint() {
    log.Printf("CONFIG: db.type=%s db.url=%s mongo.uri=%s server.addr=%s logger.level=%s",
        c.Database.Type, c.Database.URL, c.MongoDB.URI, c.Server.Address(), c.Logger.Level)
}

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

func LoadConfig(configPath ...string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if len(configPath) > 0 {
		viper.SetConfigFile(configPath[0])
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	// REMOVE PREFIX - Auto bind to normal env vars
	viper.AutomaticEnv()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		// If no config file, that's okay, use env vars
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

// LoadFromEnv loads config purely from environment variables
func LoadFromEnv() (*Config, error) {
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
	}, nil
}

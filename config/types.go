package config

type DBConfig struct {
    Driver   string `json:"driver" yaml:"driver"`
    Host     string `json:"host" yaml:"host"`
    Port     int    `json:"port" yaml:"port"`
    User     string `json:"user" yaml:"user"`
    Password string `json:"password" yaml:"password"`
    Database string `json:"database" yaml:"database"`
    SSLMode  string `json:"sslmode" yaml:"sslmode"`
}

type Config struct {
    Database DBConfig `json:"database" yaml:"database"`
}

const (
    DriverPostgres = "postgres"
    DriverMySQL    = "mysql"
    DriverMongoDB  = "mongodb"
    DriverRedis    = "redis"
)
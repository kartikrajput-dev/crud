package config

import (
	"fmt"

	"github.com/KARTIKrocks/apikit/config"
)

type Config struct {
	DBHost     string `env:"DB_HOST"     default:"localhost"`
	DBPort     int    `env:"DB_PORT"     default:"5432"`
	DBName     string `env:"DB_NAME"     default:"playground"`
	DBUser     string `env:"DB_USER"     default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" default:"postgres"`
	DBSSLMode  string `env:"DB_SSLMODE"  default:"disable"`

	ServerAddr string `env:"SERVER_ADDR" default:":8080"`
	LogLevel   string `env:"LOG_LEVEL"   default:"info"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := config.Load(&cfg, config.WithEnvFile(".env")); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode,
	)
}

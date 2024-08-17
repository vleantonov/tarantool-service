package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	StorageConfig
	ServerConfig
	AppConfig
}

// StorageConfig Storage parameters for connecting to the database
type StorageConfig struct {
	TarantoolAddress        string        `env:"TARANTOOL_ADDRESS" env-required:"true"`
	TarantoolUsername       string        `env:"TARANTOOL_USER_NAME" env-required:"true"`
	TarantoolPassword       string        `env:"TARANTOOL_PASSWORD" env-required:"true"`
	TarantoolRequestTimeout time.Duration `env:"TARANTOOL_REQUEST_TIMEOUT" env-required:"true"`
}

// ServerConfig Server parameters for running the application
type ServerConfig struct {
	Port int    `env:"APP_PORT" env-required:"true"`
	Host string `env:"APP_HOST" env-required:"true"`
}

// AppConfig Service parameters of the application for request processing
type AppConfig struct {
	TokenTTL time.Duration `env:"TOKEN_TTL" env-required:"true"`
	Secret   string        `env:"APP_SECRET_KEY" env-required:"true"`
}

// Load reads configuration from environment variables and returns a new Config instance
// or an error, if something goes wrong.
func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

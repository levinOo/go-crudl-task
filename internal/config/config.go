package config

import (
	"errors"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Конфигурация приложения
type Config struct {
	Env     string        `yaml:"env" env-default:"local"`
	Server  ServerConfig  `yaml:"server"`
	Postgre PostgreConfig `yaml:"postgre"`
}

// Конфигурация сервера
type ServerConfig struct {
	ServerPort           string        `env:"APP_PORT" env-default:":8080"`
	ServerMode           string        `yaml:"server_mode" env:"SERVER_MODE" env-default:"debug"`
	ReadTimeout          time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout         time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout          time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
	ShutdownContextValue time.Duration `yaml:"shutdown_context_value" env:"SHUTDOWN_CONTEXT_VALUE" env-default:"5s"`
}

// Конфигурация базы данных
type PostgreConfig struct {
	URL                 string        `yaml:"url" env:"POSTGRES_URL" env-required:"true"`
	PoolMax             int           `yaml:"pool_max" env:"POSTGRES_POOL_MAX" env-default:"10"`
	RetryAttempts       int           `yaml:"retry_attempts" env:"POSTGRES_RETRY_ATTEMPTS" env-default:"5"`
	RetryDelay          time.Duration `yaml:"retry_delay" env:"POSTGRES_RETRY_DELAY" env-default:"2s"`
	ContextTimeoutValue time.Duration `yaml:"context_timeout_value" env:"POSTGRES_CONTEXT_TIMEOUT_VALUE" env-default:"5s"`
}

// Загрузка конфигурации
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH не задан")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("Файл конфигурации не существует")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.New("Ошибка чтения конфигурации")
	}

	return &cfg, nil
}

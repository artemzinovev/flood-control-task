package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	RedisConfig  RedisConfig `yaml:"redis_config"`
	RequestLimit int         `yaml:"request_limit"`
}

type RedisConfig struct {
	Addr     string        `yaml:"addr"`
	Password string        `env:"REDIS_PASSWORD"`
	DB       int           `yaml:"db"`
	TTL      time.Duration `yaml:"redis_ttl"`
}

func MustLoad(log *slog.Logger) Config {
	err := godotenv.Load()
	if err != nil {
		log.Error("failed to load .env file", err)
		os.Exit(1)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Error("failed to get config_path")
		os.Exit(1)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Error("config file does not exist: %s", err)
		os.Exit(1)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Error("cannot read config: %s", err)
		os.Exit(1)
	}

	return cfg
}

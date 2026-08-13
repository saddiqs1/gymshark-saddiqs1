package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Environment Environment
	}

	Environment struct {
		AppEnv   string `env:"APP_ENV,required"`
		LogLevel string `env:"LOG_LEVEL,required"`
	}
)

func NewLocalConfig() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("no .env file found")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func loadEnv() error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		return nil
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}

	if debug {
		if err := godotenv.Load(".env"); err != nil {
			return fmt.Errorf("no .env file found")
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory")
		}

		envPath := filepath.Join(cwd, ".env")

		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("no .env file found")
		}
	}

	return nil
}

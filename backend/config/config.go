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
	err := loadEnv()
	if err != nil {
		return nil, err
	}

	return NewConfig()
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
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

		err = godotenv.Load(envPath)
		if err != nil {
			return fmt.Errorf("no .env file found")
		}
	}

	return nil
}

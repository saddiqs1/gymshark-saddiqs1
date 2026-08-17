package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Environment Environment
		Aws         Aws
		DynamoDB    DynamoDB
		AwsConfig   aws.Config
	}

	Environment struct {
		AppEnv   string `env:"APP_ENV,required"`
		LogLevel string `env:"LOG_LEVEL,required"`
	}

	Aws struct {
		Region string `env:"AWS_REGION"`
	}

	DynamoDB struct {
		Endpoint  string `env:"DYNAMODB_ENDPOINT"`
		TableName string `env:"PACK_SIZES_TABLE_NAME,required"`
	}
)

func NewConfig() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.Aws.Region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	if cfg.DynamoDB.Endpoint != "" {
		awsCfg.Credentials = credentials.NewStaticCredentialsProvider("local", "local", "")
	}
	cfg.AwsConfig = awsCfg

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

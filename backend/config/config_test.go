package config

import (
	"context"
	"testing"
)

func TestNewConfigFromEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8000")
	t.Setenv("PACK_SIZES_TABLE_NAME", "pack-sizes-test")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() returned an error: %v", err)
	}

	if cfg.Environment.AppEnv != "production" {
		t.Errorf("AppEnv = %q, want %q", cfg.Environment.AppEnv, "production")
	}

	if cfg.Environment.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.Environment.LogLevel, "debug")
	}

	if cfg.DynamoDB.Endpoint != "http://localhost:8000" {
		t.Errorf("DynamoDB Endpoint = %q, want %q", cfg.DynamoDB.Endpoint, "http://localhost:8000")
	}

	if cfg.DynamoDB.TableName != "pack-sizes-test" {
		t.Errorf("DynamoDB TableName = %q, want %q", cfg.DynamoDB.TableName, "pack-sizes-test")
	}

	credentials, err := cfg.AwsConfig.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("retrieve local AWS credentials: %v", err)
	}
	if credentials.AccessKeyID != "local" {
		t.Errorf("AWS AccessKeyID = %q, want %q", credentials.AccessKeyID, "local")
	}
}

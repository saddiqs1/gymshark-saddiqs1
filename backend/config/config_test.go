package config

import "testing"

func TestNewConfigFromEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "debug")

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
}

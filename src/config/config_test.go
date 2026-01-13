package config_test

import (
	"testing"
	"webhook-dispatcher/src/config"
)

func TestLoadConfig_DefaultWhenEnvEmpty(t *testing.T) {
	t.Setenv("LISTEN_ADDRESS", "")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddress != ":8080" {
		t.Fatalf("expected default ':8080', got '%s'", cfg.ListenAddress)
	}
}

func TestLoadConfig_WithEnv(t *testing.T) {
	t.Setenv("LISTEN_ADDRESS", "0.0.0.0:9090")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ListenAddress != "0.0.0.0:9090" {
		t.Fatalf("expected '0.0.0.0:9090', got '%s'", cfg.ListenAddress)
	}
}

package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected default DBHost localhost, got %s", cfg.DBHost)
	}
	// if cfg.MaxDBConns != 50 {
	// 	t.Errorf("expected default MaxDBConns 50, got %d", cfg.MaxDBConns)
	// }
}

func TestLoad_EnvOverrides(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("MAX_DB_CONNS", "25")
	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.MaxDBConns != 25 {
		t.Errorf("expected MaxDBConns 25, got %d", cfg.MaxDBConns)
	}
}

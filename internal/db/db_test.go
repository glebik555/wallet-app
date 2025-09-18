package db

import (
	"testing"
	"wallet-app/internal/config"
)

func TestConnect_Config(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "walletdb",
		SSLMode:    "disable",
		MaxDBConns: 5,
	}

	pool, err := Connect(cfg)

	if err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}
	if pool == nil {
		t.Fatal("expected pool to be non-nil")
	}
	defer pool.Close()
}

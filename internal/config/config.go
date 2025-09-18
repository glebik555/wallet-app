package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	MaxDBConns int
}

func Load() *Config {
	_ = godotenv.Load(".env")

	maxConns := 200
	if v := os.Getenv("MAX_DB_CONNS"); v != "" {
		if m, err := strconv.Atoi(v); err == nil {
			maxConns = m
		}
	}

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("PG_HOST", "localhost"),
		DBPort:     getEnv("PG_PORT", "5432"),
		DBUser:     getEnv("PG_USER", "postgres"),
		DBPassword: getEnv("PG_PASSWORD", "postgres"),
		DBName:     getEnv("PG_DB", "walletdb"),
		SSLMode:    getEnv("PG_SSLMODE", "disable"),
		MaxDBConns: maxConns,
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

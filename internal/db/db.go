package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"wallet-app/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode,
	)

	pgxCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("pgx parse config: %w", err)
	}

	if cfg.SSLMode != "disable" {
		if os.Getenv("PG_SSLROOTCERT") == "" {
			return nil, fmt.Errorf("sslmode=%s requires PG_SSLROOTCERT env var", cfg.SSLMode)
		}
	}

	pgxCfg.MaxConns = int32(cfg.MaxDBConns)
	pgxCfg.MaxConnLifetime = time.Hour
	pgxCfg.MinConns = 10
	pgxCfg.MaxConnIdleTime = 30 * time.Minute
	pgxCfg.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}
	return pool, nil
}

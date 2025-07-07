package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

// NewPgPool создаёт пул подключений к PostgreSQL на основе переменной окружения DATABASE_URL.
func NewPgPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL_LOCAL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL_LOCAL is not set")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("cannot connect to PostgreSQL: %w", err)
	}

	return pool, nil
}

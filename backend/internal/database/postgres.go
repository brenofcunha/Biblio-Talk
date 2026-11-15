package database

import (
        "context"
        "time"

        "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
        cfg, err := pgxpool.ParseConfig(databaseURL)
        if err != nil {
                return nil, err
        }
        cfg.MaxConns = 30
        cfg.MinConns = 5
        cfg.MaxConnLifetime = 30 * time.Minute
        cfg.MaxConnIdleTime = 5 * time.Minute
        return pgxpool.NewWithConfig(ctx, cfg)
}

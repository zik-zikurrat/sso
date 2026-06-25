package postgres

import (
	"context"
	"log/slog"
	"sso/internal/config"
	"sso/internal/helper"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres -.
type Posgtres struct {
	log  *slog.Logger
	Pool *pgxpool.Pool
}

// New -.
func New(cfg *config.Config, log *slog.Logger) (*Posgtres, error) {
	pool, err := pgxpool.New(context.Background(), helper.GetDBDsn(cfg))
	if err != nil {
		log.Error("failed to create pg pool", slog.String("error", err.Error()))
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Error("failed to ping db", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("database connected successfully")
	return &Posgtres{
		log:  log,
		Pool: pool,
	}, nil
}

// Close -.
func (p *Posgtres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

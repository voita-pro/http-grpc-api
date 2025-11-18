package repository

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/voita-pro/http-grpc-api/internal/config"
)

// Tx is a transaction abstraction to allow use with sqlc generated interfaces.
type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// DB wraps a pgx pool and gives access to Queries (sqlc) once generated.
type DB struct {
	Pool *pgxpool.Pool
}

// Connect creates a new pgx pool.
func Connect(ctx context.Context, dbCfg *config.DB) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbCfg.User,
		dbCfg.Password,
		net.JoinHostPort(dbCfg.Host, strconv.Itoa(dbCfg.Port)),
		dbCfg.Name,
	))
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 5
	cfg.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

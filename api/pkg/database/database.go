package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// Init ...
func Init( //nolint: ireturn
	host string,
	port int,
	username, password, name string,
	timeout time.Duration,
) (DBInterface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		name,
	)

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close()
		return nil, err //nolint: wrapcheck
	}

	return conn, nil
}

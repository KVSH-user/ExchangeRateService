// Package postgres contains interface
// and queries for the PostgreSQL database.
package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pgclient "github.com/KVSH-user/ExchangeRateService/internal/adapters/postgres/pgx_conn"
	"github.com/KVSH-user/ExchangeRateService/internal/config"
)

// ConnClient defines the interface for database connection operations.
type ConnClient interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	// ExecTx executes a function within a transaction, applying optional transaction parameters.
	ExecTx(ctx context.Context, f func(ctx context.Context, tx pgclient.DB) error, txParams ...pgclient.TxParams) error
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close(ctx context.Context)
}

// Store contains connections to the master and replica databases.
type Store struct {
	Master ConnClient
	logger *slog.Logger
}

// NewClient creates a new Store instance based on the provided configuration.
func NewClient(ctx context.Context, logger *slog.Logger, cfg *config.Config) (*Store, error) {
	masterConnCfg := pgclient.ClientConfig{
		Host:     fmt.Sprintf("%s:%d", cfg.Postgres.Host, cfg.Postgres.MasterPort),
		Login:    cfg.Postgres.Login,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DbName,
	}

	masterConn, err := pgclient.NewClient(ctx, cfg, &masterConnCfg, logger)
	if err != nil {
		return nil, err
	}

	return &Store{
		Master: masterConn,
		logger: logger,
	}, nil
}

// Close closes the master and replica database connections.
func (s *Store) Close(ctx context.Context) {
	s.Master.Close(ctx)
}

// 'exec' method executes an SQL command using the provided database connection and returns the command tag.
//
// Available db connection: Master, Replica, Tx (to run 'exec' inside the transaction).
func (s *Store) exec(ctx context.Context, sql string, db pgclient.DB, args ...any) (commandTag pgconn.CommandTag, err error) {
	return db.Exec(ctx, sql, args...)
}

// 'query' method executes an SQL query and returns the result set as rows using the specified database connection.
//
// Available db connection: Master, Replica, Tx (to run 'query' inside the transaction).
func (s *Store) query(ctx context.Context, sql string, db pgclient.DB, args ...any) (pgx.Rows, error) {
	return db.Query(ctx, sql, args...)
}

// 'queryRow' method executes an SQL query that is expected to return a single row, using the specified database connection.
//
// Available db connection: Master, Replica, Tx (to run 'queryRow' inside the transaction).
func (s *Store) queryRow(ctx context.Context, sql string, db pgclient.DB, args ...any) pgx.Row {
	return db.QueryRow(ctx, sql, args...)
}

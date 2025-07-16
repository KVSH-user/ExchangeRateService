// Package pgx_conn contains the implementation of the DB interface for PostgreSQL.
package pgx_conn

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"

	"github.com/KVSH-user/ExchangeRateService/internal/config"
)

const migrationsPath = "migrations/postgresql"

// DB defines the interface for database operations.
type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

// ClientConfig defines the configuration for a PostgreSQL client.
type ClientConfig struct {
	Host     string
	Login    string
	Password string
	DBName   string
}

// Client defines the PostgreSQL client.
type Client struct {
	pgxPool *pgxpool.Pool
	logger  *slog.Logger
}

// NewClient creates a new PostgreSQL client instance based on the provided configuration.
func NewClient(ctx context.Context, cfg *config.Config, pgCfg *ClientConfig, logger *slog.Logger) (*Client, error) {
	dataBaseURL := connectionString(pgCfg)
	poolCfg, err := newConfig(ctx, cfg, dataBaseURL, logger)
	if err != nil {
		return nil, err
	}

	pgxPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pgxPool)
	defer sqlDB.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("setting postgres dialect: %w", err)
	}

	if err = goose.Up(sqlDB, migrationsPath); err != nil {
		return nil, fmt.Errorf("goose up: %w", err)
	}

	if err = pgxPool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Client{
		pgxPool: pgxPool,
		logger:  logger,
	}, nil
}

// connectionString creates a PostgreSQL connection string based on the provided
// configuration.
func connectionString(cfg *ClientConfig) string {
	const baseURL = "postgresql://"
	const colonSep = ":"
	const slashSep = "/"
	const sep = "@"

	var builder strings.Builder

	if cfg.Login == "" && cfg.Password == "" {
		builder.Grow(len(baseURL) + len(cfg.Host) + len(slashSep) + len(cfg.DBName))
		builder.WriteString(baseURL)
		builder.WriteString(cfg.Host)
		builder.WriteString(slashSep)
		builder.WriteString(cfg.DBName)

		return builder.String()
	}

	builder.Grow(len(baseURL) +
		len(cfg.Login) +
		len(colonSep) +
		len(cfg.Password) +
		len(sep) +
		len(cfg.Host) +
		len(slashSep) +
		len(cfg.DBName))
	builder.WriteString(baseURL)
	builder.WriteString(cfg.Login)
	builder.WriteString(colonSep)
	builder.WriteString(cfg.Password)
	builder.WriteString(sep)
	builder.WriteString(cfg.Host)
	builder.WriteString(slashSep)
	builder.WriteString(cfg.DBName)

	return builder.String()
}

// Close closes the database connection pool.
func (c *Client) Close(_ context.Context) {
	c.pgxPool.Close()
}

// Query acquires a connection and executes a query that returns pgx.Rows.
// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
// See pgx.Rows documentation to close the returned Rows and return the acquired connection to the Pool.
//
// If there is an error, the returned pgx.Rows will be returned in an error state.
// If preferred, ignore the error returned from Query and handle errors using the returned pgx.Rows.
//
// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
// needed. See the documentation for those types for details.
func (c *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pgxPool.Query(ctx, sql, args...)
}

// Exec acquires a connection from the Pool and executes the given SQL.
// SQL can be either a prepared statement name or an SQL string.
// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
// The acquired connection is returned to the pool when the Exec function returns.
func (c *Client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.pgxPool.Exec(ctx, sql, args...)
}

// QueryRow acquires a connection and executes a query that is expected
// to return at most one row (pgx.Row). Errors are deferred until pgx.Row's
// Scan method is called. If the query selects no rows, pgx.Row's Scan will
// return ErrNoRows. Otherwise, pgx.Row's Scan scans the first selected row
// and discards the rest. The acquired connection is returned to the Pool when
// pgx.Row's Scan method is called.
//
// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
//
// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
// needed. See the documentation for those types for details.
func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pgxPool.QueryRow(ctx, sql, args...)
}

// ExecTx executes a function within a transaction, applying optional transaction parameters.
func (c *Client) ExecTx(ctx context.Context, f func(ctx context.Context, tx DB) error, txParams ...TxParams) error {
	const fName = "ExecTx"
	log := c.logger.With(fName)

	// Initialize transaction options.
	txOpts := pgx.TxOptions{}

	// Apply all transaction parameters to the options.
	for _, param := range txParams {
		txOpts = param.Apply(txOpts)
	}

	// Begin the transaction with the specified options.
	tx, err := c.pgxPool.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	// Defer rollback in case of an error during transaction execution.
	defer func(tx pgx.Tx) {
		if err == nil {
			return
		}

		// Rollback the transaction.
		if err = tx.Rollback(ctx); err != nil {
			log.InfoContext(ctx, fName, slog.String("failed to rollback transaction", err.Error()))
		}
	}(tx)

	// Execute the provided function within the transaction.
	if err = f(ctx, tx); err != nil {
		return err
	}

	// Commit the transaction if no errors occurred.
	return tx.Commit(ctx)
}

// CopyFrom uses the PostgreSQL copy protocol to perform bulk data insertion. It returns the number of rows copied and
// an error.
//
// CopyFrom requires all values use the binary format. A pgtype.Type that supports the binary format must be registered
// for the type of each column. Almost all types implemented by pgx support the binary format.
//
// Even though enum types appear to be strings they still must be registered to use with CopyFrom. This can be done with
// Conn.LoadType and pgtype.Map.RegisterType.
func (c *Client) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.pgxPool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// newConfig creates a new pgxpool.Config object based on the provided configuration.
func newConfig(ctx context.Context, cfg *config.Config, dataBaseURL string, logger *slog.Logger) (*pgxpool.Config, error) {
	const fn = "newConfig"
	log := logger.With("fn", fn)

	dbConfig, err := pgxpool.ParseConfig(dataBaseURL)
	if err != nil {
		return nil, err
	}

	setWithConfig(dbConfig, &cfg.Postgres)

	dbConfig.BeforeAcquire = func(_ context.Context, _ *pgx.Conn) bool {
		return true
	}

	dbConfig.AfterRelease = func(_ *pgx.Conn) bool {
		return true
	}

	dbConfig.BeforeClose = func(_ *pgx.Conn) {
		log.InfoContext(ctx, "Closed the connection pool to the database!")
	}

	return dbConfig, nil
}

// setWithConfig sets the configuration for the database connection pool based on
// the provided PostgreSQL configuration.
func setWithConfig(dbConfig *pgxpool.Config, pgCfg *config.PostgreSQL) {
	dbConfig.MaxConns = pgCfg.MaxConns

	dbConfig.MinConns = pgCfg.MinConns

	dbConfig.MaxConnLifetime = pgCfg.MaxConnLifetime

	dbConfig.MaxConnIdleTime = pgCfg.MaxConnIdleTime

	dbConfig.HealthCheckPeriod = pgCfg.HealthcheckPeriod

	dbConfig.ConnConfig.ConnectTimeout = pgCfg.ConnectTimeout
}

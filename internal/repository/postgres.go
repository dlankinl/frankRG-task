package repository

import (
	"context"
	"database/sql"
)

type queryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type postgresDB struct {
	db *sql.DB
	//transactor transactor // TODO: transactor
	slug     string
	hostname string
}

func NewDBConnection(db *sql.DB, slug string, hostname string) *postgresDB {
	return &postgresDB{
		db:       db,
		slug:     slug,
		hostname: hostname,
	}
}

package file

import (
	trans "FrankRGTask/pkg/transactor"
	"context"
	"database/sql"
)

type transactor interface {
	GetTx(ctx context.Context, reposSlug string) (trans.Transactor, error)
	SetTx(ctx context.Context, reposSlug string, tx trans.Transaction) error
}

type queryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

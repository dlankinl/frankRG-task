package transactor

import (
	"context"
	"errors"
)

var NoBeginOfTransactionErr = errors.New("there is no begin of transaction")
var CannotCastInterfaceErr = errors.New("cannot cast interface")
var NoTransactionErr = errors.New("no transaction for this slug")

type Transactor struct{}

func NewTransactor() Transactor {
	return Transactor{}
}

func (t Transactor) Begin(ctx context.Context) (context.Context, error) {
	transaction := map[string]Transaction{}
	return context.WithValue(ctx, "transactions", transaction), nil
}

func (t Transactor) Commit(ctx context.Context) error {
	transactions := ctx.Value("transactions")
	if transactions == nil {
		return NoBeginOfTransactionErr
	}
	transaction, status := transactions.(map[string]Transaction)
	if !status {
		return CannotCastInterfaceErr
	}
	for _, t := range transaction {
		if err := t.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (t Transactor) Rollback(ctx context.Context) error {
	transactions := ctx.Value("transactions")
	if transactions == nil {
		return NoBeginOfTransactionErr
	}
	transaction, status := transactions.(map[string]Transaction)
	if !status {
		return CannotCastInterfaceErr
	}
	for _, t := range transaction {
		if err := t.Rollback(); err != nil {
			return err
		}
	}
	return nil
}

func (t Transactor) GetTx(ctx context.Context, reposSlug string) (Transaction, error) {
	transactions := ctx.Value("transactions")
	if transactions == nil {
		return nil, NoBeginOfTransactionErr
	}
	transaction, status := transactions.(map[string]Transaction)
	if !status {
		return nil, CannotCastInterfaceErr
	}
	tx, status := transaction[reposSlug]
	if !status {
		return tx, NoTransactionErr
	}
	return tx, nil
}

func (t Transactor) SetTx(ctx context.Context, reposSlug string, tx Transaction) error {
	transactions := ctx.Value("transactions")
	if transactions == nil {
		return NoBeginOfTransactionErr
	}
	transaction, status := transactions.(map[string]Transaction)
	if !status {
		return CannotCastInterfaceErr
	}
	transaction[reposSlug] = tx
	return nil
}

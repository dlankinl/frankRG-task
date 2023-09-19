package transactor

type Transaction interface {
	Commit() error
	Rollback() error
}

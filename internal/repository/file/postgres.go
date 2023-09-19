package file

import (
	"FrankRGTask/internal/models"
	errs "FrankRGTask/pkg/errors"
	trans "FrankRGTask/pkg/transactor"
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"reflect"
)

type postgresDB struct {
	db         *sql.DB
	transactor transactor // TODO: transactor
	slug       string
	hostname   string
}

func NewDBConnection(db *sql.DB, slug string, hostname string) *postgresDB {
	return &postgresDB{
		db:       db,
		slug:     slug,
		hostname: hostname,
	}
}

func (repo *postgresDB) RegisterTX(ctx context.Context) error {
	_, err := repo.transactor.GetTx(ctx, repo.slug)
	if errors.Is(err, trans.NoTransactionErr) {
		tx, err := repo.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		err = repo.transactor.SetTx(ctx, repo.slug, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *postgresDB) getExecutor(ctx context.Context) (queryExecutor, error) {
	tx, err := repo.transactor.GetTx(ctx, repo.slug)
	if errors.Is(err, trans.NoBeginOfTransactionErr) || errors.Is(err, trans.NoTransactionErr) {
		return repo.db, nil
	}
	if err != nil {
		return nil, err
	}
	casted, status := tx.(*sql.Tx)
	if !status {
		return nil, err
	}
	return casted, nil
}

func (repo *postgresDB) Create(ctx context.Context, file models.File) error {
	executor, err := repo.getExecutor(ctx)
	if err != nil {
		return err
	}
	return create(executor, ctx, file)
}

func create(executor queryExecutor, ctx context.Context, file models.File) error {
	query := `INSERT INTO Files(name, size, modtime, isdirectory, content, parentid) 
					VALUES ($1, $2, $3, $4, $5, $6)`

	_ = executor.QueryRowContext(ctx, query, file.Name, file.Size, file.ModTime, file.IsDirectory, file.Content, file.ParentID)

	return nil
}

func (repo *postgresDB) GetParent(ctx context.Context, name string) (int, error) {
	query := `SELECT id FROM Files WHERE name = $1 AND size = 0`
	row := repo.db.QueryRowContext(ctx, query, name)
	var id int
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("entity doesn't exist")
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *postgresDB) GetContent(ctx context.Context, id int) ([]byte, error) {
	query := `SELECT Content FROM Files WHERE id = $1`
	row := repo.db.QueryRowContext(ctx, query, id)
	var content []byte
	err := row.Scan(&content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("entity doesn't exist")
	}
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (repo *postgresDB) Rename(ctx context.Context, newName string, id int) error {
	executor, err := repo.getExecutor(ctx)
	if err != nil {
		return err
	}
	return rename(executor, ctx, newName, id)
}

func rename(executor queryExecutor, ctx context.Context, newName string, id int) error {
	query := `UPDATE files
			SET name = $1
			WHERE id = $2
		`

	res, err := executor.ExecContext(ctx, query, newName, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.NothingFoundToRenameError
	}

	return nil
}

func (repo *postgresDB) FindFilesRecursive(ctx context.Context, id int) ([]int, error) {
	query := `
		WITH RECURSIVE DirectoryHierarchy AS (
		    SELECT id FROM files WHERE id = $1           
		    UNION ALL 
		    SELECT f.id FROM files f
			INNER JOIN DirectoryHierarchy dh ON f.parentid = dh.id
		)
		SELECT id FROM DirectoryHierarchy;
		`

	rows, err := repo.db.QueryContext(ctx, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.NoDirsWereFoundErr
	}
	if err != nil {
		return nil, err
	}

	var idsToDelete []int
	for rows.Next() {
		var idInner int
		if err = rows.Scan(&idInner); err != nil {
			return nil, err
		}
		idsToDelete = append(idsToDelete, idInner)
	}

	return idsToDelete, nil
}

func (repo *postgresDB) DeleteByID(ctx context.Context, id int) error {
	executor, err := repo.getExecutor(ctx)
	if err != nil {
		return err
	}
	ids, err := repo.FindFilesRecursive(ctx, id)
	if err != nil {
		return err
	}
	return deleteByID(executor, ctx, ids)
}

func deleteByID(executor queryExecutor, ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	deleteQuery := `DELETE FROM files WHERE id = ANY($1::integer[])`

	pgIntArray := pq.Array(ids)
	_, err := executor.ExecContext(ctx, deleteQuery, pgIntArray)
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgresDB) GetFilesInDir(ctx context.Context, id int) ([]models.File, error) {
	query := `SELECT * FROM Files WHERE parentid = $1`

	rows, err := repo.db.QueryContext(ctx, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.NoDirsWereFoundErr
	}
	if err != nil {
		return nil, err
	}

	var filesList []models.File
	for rows.Next() {
		var file models.File

		s := reflect.ValueOf(&file).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}
		filesList = append(filesList, file)
	}

	return filesList, err
}

func (repo *postgresDB) Content(ctx context.Context, id int) ([]byte, error) {
	query := `SELECT isdirectory, content FROM files WHERE id = $1`

	row := repo.db.QueryRowContext(ctx, query, id)
	var content []byte
	var isDirectory bool
	err := row.Scan(&isDirectory, &content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.NoFilesWereFoundErr
	}
	if isDirectory {
		return nil, errs.TypeNotFileErr
	}

	return content, nil
}

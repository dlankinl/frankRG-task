package file

import (
	"FrankRGTask/internal/models"
	errs "FrankRGTask/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"reflect"
)

type PostgresDB struct {
	db   *sql.DB
	slug string
}

func NewDBConnection(db *sql.DB, slug string) *PostgresDB {
	return &PostgresDB{
		db:   db,
		slug: slug,
	}
}

func (repo *PostgresDB) Create(ctx context.Context, file *models.File) error {
	return create(repo.db, ctx, file)
}

func create(db *sql.DB, ctx context.Context, file *models.File) error {
	query := `INSERT INTO Files(name, size, modtime, isdirectory, content, parentid) 
					VALUES ($1, $2, $3, $4, $5, $6)`

	err := db.QueryRowContext(ctx, query, file.Name, file.Size, file.ModTime, file.IsDirectory, file.Content, file.ParentID)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (repo *PostgresDB) GetParent(ctx context.Context, name string) (int, error) {
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

func (repo *PostgresDB) Rename(ctx context.Context, newName string, id int) error {
	return rename(repo.db, ctx, newName, id)
}

func rename(db *sql.DB, ctx context.Context, newName string, id int) error {
	query := `UPDATE files
			SET name = $1
			WHERE id = $2
		`

	res, err := db.ExecContext(ctx, query, newName, id)
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

func (repo *PostgresDB) FindFilesRecursive(ctx context.Context, id int) ([]int, error) {
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

func (repo *PostgresDB) DeleteByID(ctx context.Context, id int) (int, error) {
	ids, err := repo.FindFilesRecursive(ctx, id)
	if err != nil {
		return 0, err
	}
	return deleteByID(repo.db, ctx, ids)
}

func deleteByID(db *sql.DB, ctx context.Context, ids []int) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	deleteQuery := `DELETE FROM files WHERE id = ANY($1::integer[])`

	pgIntArray := pq.Array(ids)
	_, err := db.ExecContext(ctx, deleteQuery, pgIntArray)
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}

func (repo *PostgresDB) GetFilesInDir(ctx context.Context, id int) ([]models.File, error) {
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

func (repo *PostgresDB) GetContent(ctx context.Context, id int) ([]byte, error) {
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

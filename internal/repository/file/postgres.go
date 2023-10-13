package file

import (
	errs "FrankRGTask/internal/errors"
	"FrankRGTask/internal/models"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"io"
	"reflect"
)

type PostgresDB struct {
	db *pgx.Conn
}

func NewDBConnection(db *pgx.Conn) PostgresDB {
	return PostgresDB{
		db: db,
	}
}

func (repo *PostgresDB) Create(ctx context.Context, file *models.File, fileReader io.Reader) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx (create): %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %w", rollbackErr)
			}
		}
	}()

	loStorage := tx.LargeObjects()

	loId, err := loStorage.Create(ctx, 0)
	if err != nil {
		return fmt.Errorf("creating large object: %w", err)
	}

	var query string
	if fileReader != nil {
		lo, err := loStorage.Open(ctx, loId, pgx.LargeObjectModeWrite)
		if err != nil {
			return fmt.Errorf("opening large object: %w", err)
		}

		hash := sha256.New()
		teeReader := io.TeeReader(fileReader, hash)

		_, err = io.Copy(lo, teeReader)
		if err != nil {
			return fmt.Errorf("copying data to db large object from file: %w", err)
		}
	}

	query = `
			INSERT INTO File_Data(hash, data_oid)
			VALUES ($1, $2)
			RETURNING id
		`

	var dataId int
	err = repo.db.QueryRow(ctx, query, base64.URLEncoding.EncodeToString(file.Content), loId).Scan(&dataId)
	if err != nil {
		return fmt.Errorf("inserting file data: %w", err)
	}

	query = `
		INSERT INTO Files(name, size, modtime, isdirectory, parentid, file_data_id) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = repo.db.Exec(ctx, query, file.Name, file.Size, file.ModTime, file.IsDirectory, file.ParentID, dataId)
	if err != nil {
		return fmt.Errorf("creating entity: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (repo *PostgresDB) GetParent(ctx context.Context, name string) (int, error) {
	query := `SELECT id FROM Files 
          		WHERE name = $1 AND size = 0 AND isdirectory = true`
	row := repo.db.QueryRow(ctx, query, name)
	var id int
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("entity doesn't exist")
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PostgresDB) Rename(ctx context.Context, newName string, id int) error {
	query := `UPDATE files
			SET name = $1
			WHERE id = $2
		`

	res, err := repo.db.Exec(ctx, query, newName, id)
	if err != nil {
		return fmt.Errorf("renaming entity: %w", err)
	}

	rowsAffected := res.RowsAffected()

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

	rows, err := repo.db.Query(ctx, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.NoDirsWereFoundErr
	}
	if err != nil {
		return nil, fmt.Errorf("recursive searching for files (query): %w", err)
	}

	var idsToDelete []int
	for rows.Next() {
		var idInner int
		if err = rows.Scan(&idInner); err != nil {
			return nil, fmt.Errorf("searching for files (scan): %w", err)
		}
		idsToDelete = append(idsToDelete, idInner)
	}

	return idsToDelete, nil
}

func (repo *PostgresDB) DeleteByID(ctx context.Context, id int) (int, error) {
	ids, err := repo.FindFilesRecursive(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("deleting by ID (recursive searching for files): %w", err)
	}

	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx (create): %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %w", rollbackErr)
			}
		}
	}()

	deleteQuery := `
		DELETE FROM file_data 
		       WHERE id IN (
		           SELECT file_data_id FROM files
		        	WHERE id = ANY($1::integer[])
		)
	`

	pgIntArray := pq.Array(ids)
	_, err = repo.db.Exec(ctx, deleteQuery, pgIntArray)
	if err != nil {
		return 0, fmt.Errorf("deleting files data (exec): %w", err)
	}

	deleteQuery = `DELETE FROM files WHERE id = ANY($1::integer[])`

	_, err = repo.db.Exec(ctx, deleteQuery, pgIntArray)
	if err != nil {
		return 0, fmt.Errorf("deleting files by IDs (exec): %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return len(ids), nil
}

func (repo *PostgresDB) GetFilesInDir(ctx context.Context, id int) ([]models.File, error) {
	query := `SELECT * FROM Files WHERE parentid = $1`

	rows, err := repo.db.Query(ctx, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.NoDirsWereFoundErr
	}
	if err != nil {
		return nil, fmt.Errorf("getting files in dir (query): %w", err)
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
			return nil, fmt.Errorf("getting files in dir (scan): %w", err)
		}
		filesList = append(filesList, file)
	}

	return filesList, err
}

func (repo *PostgresDB) GetContent(ctx context.Context, id int) ([]byte, error) {
	type FileData struct {
		IsDirectory bool
		Content     []byte
	}

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %w", rollbackErr)
			}
		}
	}()

	query := `
		SELECT fd.data_oid, f.isdirectory, fd.hash 
		FROM files f 
		INNER JOIN file_data fd
			ON fd.id = f.file_data_id
		WHERE fd.id = $1
	`

	var loId uint32
	var FData FileData
	var tmp string
	err = repo.db.QueryRow(ctx, query, id).Scan(
		&loId,
		&FData.IsDirectory,
		&tmp,
	)

	FData.Content, err = base64.URLEncoding.DecodeString(tmp)
	if err != nil {
		return nil, fmt.Errorf("content decode: %w", err)
	}

	return FData.Content, nil
}

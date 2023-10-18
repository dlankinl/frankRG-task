package file

import (
	errs "FrankRGTask/internal/errors"
	"FrankRGTask/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func (repo PostgresDB) CreateDir(ctx context.Context, file *models.File) error {
	var query string
	query = `
		insert into files(name, size, mod_time, is_directory, parent_id) 
		values (@name, @size, @mod_time, @is_directory, @parent_id)
	`

	_, err := repo.db.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"name":         file.Name,
			"size":         file.Size,
			"mod_time":     file.ModTime,
			"is_directory": file.IsDirectory,
			"parent_id":    file.ParentID,
		},
	)

	if err != nil {
		return fmt.Errorf("creating dir: %w", err)
	}

	return nil
}

func (repo PostgresDB) Create(ctx context.Context, file *models.File, content []byte) error {
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

	lo, err := loStorage.Open(ctx, loId, pgx.LargeObjectModeWrite)
	if err != nil {
		return fmt.Errorf("opening large object: %w", err)
	}
	_, err = lo.Write(content)
	if err != nil {
		return fmt.Errorf("writing data to db large object from request: %w")
	}

	var query string
	query = `
		insert into files(name, size, mod_time, is_directory, parent_id, data_oid) 
		values (@name, @size, @mod_time, @is_directory, @parent_id, @data_oid)
	`

	_, err = tx.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"name":         file.Name,
			"size":         file.Size,
			"mod_time":     file.ModTime,
			"is_directory": file.IsDirectory,
			"parent_id":    file.ParentID,
			"data_oid":     loId,
		},
	)

	if err != nil {
		return fmt.Errorf("creating entity: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (repo PostgresDB) Upload(ctx context.Context, file *models.File, reader io.Reader) error {
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

	lo, err := loStorage.Open(ctx, loId, pgx.LargeObjectModeWrite)
	if err != nil {
		return fmt.Errorf("opening large object: %w", err)
	}

	_, err = io.Copy(lo, reader)
	if err != nil {
		return fmt.Errorf("copying data to db large object from file: %w", err)
	}

	var query string
	query = `
		insert into files(name, size, mod_time, is_directory, parent_id, data_oid) 
		values (@name, @size, @mod_time, @is_directory, @parent_id, @data_oid)
	`

	_, err = tx.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"name":         file.Name,
			"size":         file.Size,
			"mod_time":     file.ModTime,
			"is_directory": file.IsDirectory,
			"parent_id":    file.ParentID,
			"data_oid":     loId,
		},
	)

	if err != nil {
		return fmt.Errorf("creating entity: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (repo PostgresDB) GetParent(ctx context.Context, name string) (int, error) {
	query := `select id 
				from files 
          		where name = @name and size = 0 and is_directory = true`
	row := repo.db.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"name": name,
		},
	)
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

func (repo PostgresDB) Rename(ctx context.Context, newName string, id int) error {
	query := `update files
			set name = @name
			where id = @id
		`

	res, err := repo.db.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"name": newName,
			"id":   id,
		},
	)
	if err != nil {
		return fmt.Errorf("renaming entity: %w", err)
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		return errs.NothingFoundToRenameError
	}

	return nil
}

func (repo PostgresDB) FindFilesRecursive(ctx context.Context, id int) ([]int, error) {
	query := `
		with recursive DirectoryHierarchy as (
		    select id from files where id = @id           
		    union all 
		    select f.id from files f
			inner join DirectoryHierarchy dh on f.parent_id = dh.id
		)
		select id 
		from DirectoryHierarchy;
	`

	rows, err := repo.db.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"id": id,
		},
	)

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

func (repo PostgresDB) GetFilesInDir(ctx context.Context, id int) ([]models.File, error) {
	query := `
		select * 
		from files 
		where parent_id = @parent_id
	`

	rows, err := repo.db.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"parent_id": id,
		},
	)

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

func (repo PostgresDB) Download(ctx context.Context, id int, readFn func(reader io.Reader) error) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if err != nil {
				err = fmt.Errorf("rollback err: %w", rollbackErr)
			}
		}
	}()

	var loId uint32
	query := `
		select f.data_oid
		from files f
		where f.id = @id
	`

	err = repo.db.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"id": id,
		},
	).Scan(&loId)

	loStorage := tx.LargeObjects()

	lo, err := loStorage.Open(ctx, loId, pgx.LargeObjectModeRead)
	if err != nil {
		return fmt.Errorf("opening large files with id=%d: %w", loId, err)
	}

	err = readFn(lo)
	if err != nil {
		return fmt.Errorf("copying large object data to file: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}

	return nil
}

func (repo PostgresDB) DeleteFile(ctx context.Context, id int) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %w", rollbackErr)
			}
		}
	}()

	ids, err := repo.FindFilesRecursive(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting by ID (recursive searching for files): %w", err)
	}

	if len(ids) == 0 {
		return fmt.Errorf("delete files: %w", err)
	}

	query := `
		select f.data_oid
		from files f
		where f.id = any(@ids::int[]) 
		  and f.is_directory = false;
	`

	rows, err := tx.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"ids": ids,
		},
	)

	var oids []uint32
	for rows.Next() {
		var oid uint32

		err = rows.Scan(&oid)
		if err != nil {
			return fmt.Errorf("getting oid (scan): %w", err)
		}
		oids = append(oids, oid)
	}

	for _, oid := range oids {
		_, err = tx.Exec(
			ctx,
			`select lo_unlink(@oid)`,
			pgx.NamedArgs{
				"oid": oid,
			},
		)
		if err != nil {
			return fmt.Errorf("unlink large object: %w", err)
		}
	}

	query = `
		delete from files
		where id = any(@ids::int[])
	`

	_, err = tx.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"ids": ids,
		},
	)

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

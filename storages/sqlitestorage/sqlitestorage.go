package sqlitestorage

import (
	"database/sql"
	"fmt"
	"lena/errors"
	"lena/storages"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
	storages.Storage
}

func NewSqliteStorage() (*SqliteStorage, error) {
	domain := "sqlitestorage.NewSqliteStorage"
	db, err := sql.Open("sqlite3", "lena_storage.sqlite")
	if err != nil {
		return nil, errors.NewAppError(errors.ErrCodeOpeningSqliteDB, domain, err)
	}
	err = migrate(db)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrCodeMigratingSqliteDB, domain, err)
	}
	return &SqliteStorage{db: db}, nil
}

func migrate(db *sql.DB) error {
	domain := "sqlitestorage.migrate"
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS migrations (
				version INTEGER PRIMARY KEY
			);
		`,
	)
	if err != nil {
		return errors.NewAppError(errors.ErrCodeCreatingMigrationsTable, domain, err)
	}
	rows, err := db.Query("SELECT version FROM migrations")
	if err != nil {
		return errors.NewAppError(errors.ErrCodeQueryingMigrations, domain, err)
	}
	defer rows.Close()
	appliedMigration := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return errors.NewAppError(errors.ErrCodeRowScanHasFailed, domain, err)
		}
		appliedMigration[version] = true
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.NewAppError(errors.ErrCodeCannotBeginDBTx, domain, err)
	}
	defer tx.Rollback()
	for _, migration := range migrations() {
		if appliedMigration[migration.version] {
			continue
		}
		domainWithVersion := fmt.Sprintf("%s: migration version = %d", domain, migration.version)
		_, err = tx.Exec(migration.command)
		if err != nil {
			return errors.NewAppError(errors.ErrCodeExecutingMigrationCmd, domainWithVersion, err)
		}
		err = migration.modify(tx)
		if err != nil {
			return errors.NewAppError(errors.ErrCodeMigrationModify, domainWithVersion, err)
		}
		_, err = tx.Exec("INSERT INTO migrations (version) VALUES (?)", migration.version)
		if err != nil {
			return errors.NewAppError(errors.ErrCodeInsertingMigration, domainWithVersion, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.NewAppError(errors.ErrCodeDBTxCommitHasFailed, domain, err)
	}
	return nil
}

func toTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05-07:00", str)
}

func toBool(i int) bool {
	return i != 0
}

package sqlitestorage

import (
	"database/sql"
	"lena/storages"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
	storages.Storage
}

func NewSqliteStorage() (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", "storage.sqlite")
	if err != nil {
		return nil, err
	}
	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return &SqliteStorage{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS migrations (
				version INTEGER PRIMARY KEY
			);
		`,
	)
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT version FROM migrations")
	if err != nil {
		return err
	}
	defer rows.Close()
	appliedMigration := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return err
		}
		appliedMigration[version] = true
	}
	for _, migration := range migrations() {
		if appliedMigration[migration.version] {
			continue
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		_, err = tx.Exec(migration.command)
		if err != nil {
			return err
		}
		err = migration.modify(tx)
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO migrations (version) VALUES (?)", migration.version)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func toTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05-07:00", str)
}

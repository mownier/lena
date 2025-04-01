package sqlitestorage

import (
	"context"
	"database/sql"
	"errors"
	"lena/models"
)

func (s *SqliteStorage) AddUser(ctx context.Context, user models.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	exists, err := s.checkUserIfExistingByName(ctx, user.Name, tx)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	_, err = tx.ExecContext(ctx,
		`
		INSERT INTO users 
			(name, password, created_on)
			VALUES(?, ?, ?)
		`,
		user.Name, user.Password, user.CreatedOn,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteStorage) GetUserByName(ctx context.Context, name string) (models.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx,
		`
		SELECT name, password, created_on
			FROM users
			WHERE name = ?
		`,
		name,
	)
	var createdOn string
	user := models.User{}
	err = row.Scan(&user.Name, &user.Password, &createdOn)
	if err != nil {
		return models.User{}, err
	}
	user.CreatedOn, err = toTime(createdOn)
	if err != nil {
		return models.User{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *SqliteStorage) checkUserIfExistingByName(ctx context.Context, name string, tx *sql.Tx) (bool, error) {
	row := tx.QueryRowContext(ctx,
		`
		SELECT 1
			FROM users
			WHERE name = ?
			LIMIT 1
		`,
		name,
	)
	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

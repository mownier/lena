package sqlitestorage

import (
	"context"
	"database/sql"
	"fmt"
	"lena/errors"
	"lena/models"
)

func (s *SqliteStorage) AddUser(ctx context.Context, user models.User) error {
	domain := fmt.Sprintf("sqlitestorage.SqliteStorage.AddUser: user = %v", user)
	tx, err := s.db.Begin()
	if err != nil {
		return errors.NewAppError(errors.ErrCodeCannotBeginDBTx, domain, err)
	}
	defer tx.Rollback()
	exists, err := s.checkUserIfExistingByName(ctx, user.Name, tx)
	if err != nil {
		return errors.NewAppError(errors.ErrCodeCheckingIfUserExists, domain, err)
	}
	if exists {
		return errors.NewAppError(errors.ErrCodeUserAlreadyExists, domain, nil)
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
		return errors.NewAppError(errors.ErrCodeQueryingUsers, domain, err)
	}
	err = tx.Commit()
	if err != nil {
		return errors.NewAppError(errors.ErrCodeDBTxCommitHasFailed, domain, err)
	}
	return nil
}

func (s *SqliteStorage) GetUserByName(ctx context.Context, name string) (models.User, error) {
	domain := fmt.Sprintf("sqlitestorage.SqliteStorage.GetUserByName: name = %s", name)
	tx, err := s.db.Begin()
	if err != nil {
		return models.User{}, errors.NewAppError(errors.ErrCodeCannotBeginDBTx, domain, err)
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
		return models.User{}, errors.NewAppError(errors.ErrCodeRowScanHasFailed, domain, err)
	}
	user.CreatedOn, err = toTime(createdOn)
	if err != nil {
		return models.User{}, errors.NewAppError(errors.ErrCodeUserCreationTimeCannotBeDetermined, domain, err)
	}
	err = tx.Commit()
	if err != nil {
		return models.User{}, errors.NewAppError(errors.ErrCodeDBTxCommitHasFailed, domain, err)
	}
	return user, nil
}

func (s *SqliteStorage) checkUserIfExistingByName(ctx context.Context, name string, tx *sql.Tx) (bool, error) {
	domain := fmt.Sprintf("sqlitestorage.SqliteStorage.checkUserIfExistingByName: name = %s", name)
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
		return false, errors.NewAppError(errors.ErrCodeRowScanHasFailed, domain, err)
	}
	return true, nil
}

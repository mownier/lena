package sqlitestorage

import (
	"context"
	"database/sql"
	"errors"
	"lena/models"
	"lena/util"
)

func (s *SqliteStorage) AddUser(ctx context.Context, user models.User) error {
	exists, err := s.checkUserIfExistingByName(ctx, user.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	_, err = s.db.ExecContext(ctx,
		`
		INSERT INTO users 
			(name, password, created_on)
			VALUES(?, ?, ?)
		`,
		user.Name, user.Password, float64(user.CreatedOn.Unix()),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteStorage) GetUserByName(ctx context.Context, name string) (models.User, error) {
	row := s.db.QueryRowContext(ctx,
		`
		SELECT name, password, created_on
			FROM users
			WHERE name = ?
		`,
		name,
	)
	user := models.User{}
	var createdOn float64
	err := row.Scan(&user.Name, &user.Password, &createdOn)
	if err != nil {
		return models.User{}, err
	}
	user.CreatedOn = util.FromFloat64ToTime(createdOn).UTC()
	return user, nil
}

func (s *SqliteStorage) checkUserIfExistingByName(ctx context.Context, name string) (bool, error) {
	row := s.db.QueryRowContext(ctx,
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

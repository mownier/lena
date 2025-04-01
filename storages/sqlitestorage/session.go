package sqlitestorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lena/models"
	"lena/util"
	"strings"
)

func (s *SqliteStorage) AddSession(ctx context.Context, session models.Session) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	row := s.db.QueryRowContext(ctx,
		`
		SELECT 1
			FROM sessions
			WHERE access_token = ?
			LIMIT 1
		`,
		session.AccessToken,
	)
	var count int
	err = row.Scan(&count)
	var exists bool
	if err == sql.ErrNoRows {
		exists = false
	} else if err != nil {
		return err
	} else {
		exists = true
	}
	if exists {
		return errors.New("session already exists")
	}
	_, err = s.db.ExecContext(ctx,
		`
		INSERT INTO sessions
			(access_token, refresh_token, user_name, 
				access_token_expiry, refresh_token_expiry, created_on,
				archived_on, archived)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		session.AccessToken, session.RefreshToken, session.UserName,
		float64(session.AccesTokenExpiry.Unix()), float64(session.RefreshTokenExpiry.Unix()), float64(session.CreatedOn.Unix()),
		float64(session.ArchivedOn.Unix()), util.FromBoolToInt(session.Archived),
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

func (s *SqliteStorage) GetSessionByAccessToken(ctx context.Context, accessToken string) (models.Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return models.Session{}, err
	}
	defer tx.Rollback()
	session, err := s.getSessionByAccessToken(ctx, accessToken, tx)
	if err != nil {
		return models.Session{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (s *SqliteStorage) UpdateSessionByAccessToken(ctx context.Context, accessToken string, update models.SessionUpdate) (models.Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return models.Session{}, err
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx,
		`
		SELECT 1
			FROM sessions
			WHERE access_token = ?
			LIMIT 1
		`,
		accessToken,
	)
	var count int
	err = row.Scan(&count)
	if err == sql.ErrNoRows {
		return models.Session{}, errors.New("session does not exist")
	}
	if err != nil {
		return models.Session{}, err
	}
	hasUpdate := false
	columnNames := []string{}
	columnValues := []any{}
	if update.ArchivedOn != nil {
		columnNames = append(columnNames, "archived_on = ?")
		columnValues = append(columnValues, util.FromTimeToFloat64(*update.ArchivedOn))
		hasUpdate = true
	}
	if update.Archived != nil {
		columnNames = append(columnNames, "archived = ?")
		columnValues = append(columnValues, util.FromBoolToInt(*update.Archived))
		hasUpdate = true
	}
	if !hasUpdate {
		return models.Session{}, errors.New("no session to update")
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf(
			`
			UPDATE sessions
			SET %v
			WHERE access_token = ?
			`, strings.Join(columnNames, ", "),
		),
	)
	if err != nil {
		return models.Session{}, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(append(columnValues, accessToken)...)
	if err != nil {
		return models.Session{}, err
	}
	session, err := s.getSessionByAccessToken(ctx, accessToken, tx)
	if err != nil {
		return models.Session{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (s *SqliteStorage) getSessionByAccessToken(ctx context.Context, accessToken string, tx *sql.Tx) (models.Session, error) {
	row := tx.QueryRowContext(ctx,
		`
		SELECT access_token, refresh_token, user_name, 
				access_token_expiry, refresh_token_expiry, created_on,
				archived_on, archived
			FROM sessions 
			WHERE access_token = ?
		`,
		accessToken,
	)
	var session models.Session
	var accessTokenExpiry float64
	var refresTokenExpiry float64
	var createdOn float64
	var archivedOn float64
	var archived int
	err := row.Scan(
		&session.AccessToken, &session.RefreshToken, &session.UserName,
		&accessTokenExpiry, &refresTokenExpiry, &createdOn,
		&archivedOn, &archived,
	)
	if err != nil {
		return models.Session{}, err
	}
	session.AccesTokenExpiry = util.FromFloat64ToTime(accessTokenExpiry).UTC()
	session.RefreshTokenExpiry = util.FromFloat64ToTime(refresTokenExpiry).UTC()
	session.CreatedOn = util.FromFloat64ToTime(createdOn).UTC()
	session.ArchivedOn = util.FromFloat64ToTime(archivedOn).UTC()
	session.Archived = util.FromIntToBool(archived)
	return session, nil
}

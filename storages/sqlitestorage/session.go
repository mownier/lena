package sqlitestorage

import (
	"context"
	"database/sql"
	"fmt"
	"lena/errors"
	"lena/models"
	"strings"
)

func (s *SqliteStorage) AddSession(ctx context.Context, session models.Session) error {
	domain := fmt.Sprintf("sqliteStorage.SqliteStorage.AddSession: session = %v", session)
	tx, err := s.db.Begin()
	if err != nil {
		return errors.NewAppError(errors.ErrCodeCannotBeginDBTx, domain, err)
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx,
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
		return errors.NewAppError(errors.ErrCodeRowScanHasFailed, domain, nil)
	} else {
		exists = true
	}
	if exists {
		return errors.NewAppError(errors.ErrCodeSessionAlreadyExists, domain, nil)
	}
	_, err = tx.ExecContext(ctx,
		`
		INSERT INTO sessions
			(access_token, refresh_token, user_name, 
				access_token_expiry, refresh_token_expiry, created_on,
				archived_on, archived)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		session.AccessToken, session.RefreshToken, session.UserName,
		session.AccesTokenExpiry, session.RefreshTokenExpiry, session.CreatedOn,
		session.ArchivedOn, session.Archived,
	)
	if err != nil {
		return errors.NewAppError(errors.ErrCodeInsertingSession, domain, err)
	}
	err = tx.Commit()
	if err != nil {
		return errors.NewAppError(errors.ErrCodeDBTxCommitHasFailed, domain, err)
	}
	return nil
}

func (s *SqliteStorage) GetSessionByAccessToken(ctx context.Context, accessToken string) (models.Session, error) {
	domain := fmt.Sprintf("sqliteStorage.SqliteStorage.GetSessionByAccessToken: accessToken = %s", accessToken)
	tx, err := s.db.Begin()
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeCannotBeginDBTx, domain, err)
	}
	defer tx.Rollback()
	session, err := s.getSessionByAccessToken(ctx, accessToken, tx)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeGettingSession, domain, err)
	}
	err = tx.Commit()
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeDBTxCommitHasFailed, domain, err)
	}
	return session, nil
}

func (s *SqliteStorage) UpdateSessionByAccessToken(ctx context.Context, accessToken string, update models.SessionUpdate) (models.Session, error) {
	domain := fmt.Sprintf("sqliteStorage.SqliteStorage.UpdateSessionByAccessToken: accessToken = %s, update = %v", accessToken, update)
	tx, err := s.db.Begin()
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeCannotBeginDBTx, domain, err)
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
		return models.Session{}, errors.NewAppError(errors.ErrCodeSessionDoesNotExist, domain, nil)
	}
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeRowScanHasFailed, domain, err)
	}
	hasUpdate := false
	columnNames := []string{}
	columnValues := []any{}
	if update.ArchivedOn != nil {
		columnNames = append(columnNames, "archived_on = ?")
		columnValues = append(columnValues, *update.ArchivedOn)
		hasUpdate = true
	}
	if update.Archived != nil {
		columnNames = append(columnNames, "archived = ?")
		columnValues = append(columnValues, *update.Archived)
		hasUpdate = true
	}
	if !hasUpdate {
		return models.Session{}, errors.NewAppError(errors.ErrCodeNoSessionToUpdate, domain, nil)
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
		return models.Session{}, errors.NewAppError(errors.ErrCodePreparingSessionUpdateStmt, domain, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(append(columnValues, accessToken)...)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeExecutingSessionUpdateStmt, domain, err)
	}
	session, err := s.getSessionByAccessToken(ctx, accessToken, tx)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeGettingSession, domain, err)
	}
	err = tx.Commit()
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeDBTxCommitHasFailed, domain, err)
	}
	return session, nil
}

func (s *SqliteStorage) getSessionByAccessToken(ctx context.Context, accessToken string, tx *sql.Tx) (models.Session, error) {
	domain := fmt.Sprintf("sqliteStorage.SqliteStorage.getSessionByAccessToken: accessToken = %s", accessToken)
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
	var accessTokenExpiry string
	var refresTokenExpiry string
	var createdOn string
	var archivedOn string
	var archived int
	err := row.Scan(
		&session.AccessToken, &session.RefreshToken, &session.UserName,
		&accessTokenExpiry, &refresTokenExpiry, &createdOn,
		&archivedOn, &archived,
	)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeRowScanHasFailed, domain, err)
	}
	session.AccesTokenExpiry, err = toTime(accessTokenExpiry)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeSessionAcccessTokenExpiryCannotBeDetermined, domain, err)
	}
	session.RefreshTokenExpiry, err = toTime(refresTokenExpiry)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeSessionRefreshTokenExpiryCannotBeDetermined, domain, err)
	}
	session.CreatedOn, err = toTime(createdOn)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeSessionCreationTimeCannotBeDetermined, domain, err)
	}
	session.ArchivedOn, err = toTime(archivedOn)
	if err != nil {
		return models.Session{}, errors.NewAppError(errors.ErrCodeSessionArchivedTimeCannotBeDetermined, domain, err)
	}
	session.Archived = toBool(archived)
	return session, nil
}

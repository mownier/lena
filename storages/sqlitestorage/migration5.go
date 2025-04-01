package sqlitestorage

import (
	"database/sql"
	"errors"
	"time"
)

func migration5() migration {
	return migration{
		version: 5,
		command: `
			ALTER TABLE sessions ADD COLUMN created_on_2 TEXT NOT NULL DEFAULT '';
			ALTER TABLE sessions ADD COLUMN archived_on_2 TEXT NOT NULL DEFAULT '';
			ALTER TABLE sessions ADD COLUMN access_token_expiry_2 TEXT NOT NULL DEFAULT '';
			ALTER TABLE sessions ADD COLUMN refresh_token_expiry_2 TEXT NOT NULL DEFAULT '';
		`,
		modify: func(tx *sql.Tx) error {
			rows, err := tx.Query(`
				SELECT access_token, 
						created_on, archived_on, 
						access_token_expiry, refresh_token_expiry
					FROM sessions
				`,
			)
			if err != nil {
				return err
			}
			for rows.Next() {
				var accessToken string
				var createdOn float64
				var archivedOn any
				var accessTokenExpiry float64
				var refreshTokenExpiry float64
				if err = rows.Scan(&accessToken,
					&createdOn, &archivedOn,
					&accessTokenExpiry, &refreshTokenExpiry,
				); err != nil {
					return err
				}
				stmt, err := tx.Prepare(
					`
					UPDATE sessions
						SET created_on_2 = ?, archived_on_2 = ?,
							access_token_expiry_2 = ?, refresh_token_expiry_2 = ?
						WHERE access_token = ?
					`,
				)
				if err != nil {
					return err
				}
				defer stmt.Close()
				floatAsTime := func(v float64) time.Time {
					seconds := int64(v)
					nanoseconds := int64((v - float64(v)) * 1e9)
					return time.Unix(seconds, nanoseconds)
				}
				anyAsTime := func(a any) (time.Time, error) {
					switch v := a.(type) {
					case float64:
						return floatAsTime(float64(v)), nil
					case string:
						return toTime(v)
					default:
						return time.Now(), errors.New("any is neither float64 nor string")
					}
				}
				newArchivedOn, err := anyAsTime(archivedOn)
				if err != nil {
					return err
				}
				_, err = stmt.Exec(
					floatAsTime(createdOn), newArchivedOn,
					floatAsTime(accessTokenExpiry), floatAsTime(refreshTokenExpiry),
					accessToken,
				)
				if err != nil {
					return err
				}
			}
			_, err = tx.Exec(
				`
				ALTER TABLE sessions DROP COLUMN created_on;
				ALTER TABLE sessions DROP COLUMN archived_on;
				ALTER TABLE sessions DROP COLUMN access_token_expiry;
				ALTER TABLE sessions DROP COLUMN refresh_token_expiry;
				ALTER TABLE sessions RENAME created_on_2 TO created_on;
				ALTER TABLE sessions RENAME archived_on_2 TO archived_on;
				ALTER TABLE sessions RENAME access_token_expiry_2 TO access_token_expiry;
				ALTER TABLE sessions RENAME refresh_token_expiry_2 TO refresh_token_expiry;
				`,
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

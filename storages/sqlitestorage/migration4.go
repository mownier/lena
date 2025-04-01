package sqlitestorage

import (
	"database/sql"
	"time"
)

func migration4() migration {
	return migration{
		version: 4,
		command: `
			ALTER TABLE users 
				ADD COLUMN created_on_2 TEXT NOT NULL DEFAULT ''
			`,
		modify: func(tx *sql.Tx) error {
			rows, err := tx.Query("SELECT name, created_on FROM users")
			if err != nil {
				return err
			}
			users := make(map[string]float64)
			for rows.Next() {
				var name string
				var createdOn float64
				if err := rows.Scan(&name, &createdOn); err != nil {
					return err
				}
				users[name] = createdOn
			}
			for name, createdOn := range users {
				stmt, err := tx.Prepare(
					`
					UPDATE users
						SET created_on_2 = ?
						WHERE name = ?
					`,
				)
				if err != nil {
					return err
				}
				defer stmt.Close()
				seconds := int64(createdOn)
				nanoseconds := int64((createdOn - float64(createdOn)) * 1e9)
				newCreatedOn := time.Unix(seconds, nanoseconds)
				_, err = stmt.Exec(newCreatedOn, name)
				if err != nil {
					return err
				}
			}
			_, err = tx.Exec(`
				ALTER TABLE users DROP COLUMN created_on;
				ALTER TABLE users RENAME created_on_2 TO created_on;
				`,
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

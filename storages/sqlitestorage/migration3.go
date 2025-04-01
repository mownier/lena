package sqlitestorage

func migration3() migration {
	return migration{
		version: 3,
		command: `
			CREATE TABLE IF NOT EXISTS sessions (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				access_token TEXT UNIQUE NOT NULL,
				refresh_token TEXT NOT NULL,
				user_name TEXT NOT NULL,
				access_token_expiry REAL NOT NULL,
				refresh_token_expiry REAL NOT NULL,
				archived_on REAL NOT NULL,
				archived INTEGER NOT NULL,
				created_on REAL NOT NULL,

				FOREIGN KEY (user_name) REFERENCES users(name)
			);
		`,
		modify: migrationDefaultModify,
	}
}

package sqlitestorage

func migration1() migration {
	return migration{
		version: 1,
		command: `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT UNIQUE NOT NULL,
				created_on REAL NOT NULL
			);
		`,
		modify: migrationDefaultModify,
	}
}

package sqlitestorage

func migration2() migration {
	return migration{
		version: 2,
		command: `
			ALTER TABLE users ADD COLUMN password TEXT NOT NULL
		`,
		modify: migrationDefaultModify,
	}
}

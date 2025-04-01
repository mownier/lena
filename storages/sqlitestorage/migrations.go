package sqlitestorage

import (
	"database/sql"
	"sort"
)

type migration struct {
	version int
	command string
	modify  func(tx *sql.Tx) error
}

func migrationDefaultModify(tx *sql.Tx) error {
	return nil
}

func migrations() []migration {
	list := []migration{
		migration1(),
		migration2(),
		migration3(),
		migration4(),
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].version < list[j].version
	})
	return list
}

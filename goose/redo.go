package goose

import (
	"database/sql"
)

// Redo rolls back the most recently applied migration, then runs it again.
func Redo(db *sql.DB, dir string) (err error) {
	var currentVersion int64
	var migrations Migrations
	var current *Migration

	if currentVersion, err = GetDBVersion(db); err != nil {
		return
	}
	if migrations, err = CollectMigrations(dir, minVersion, maxVersion); err != nil {
		return
	}
	if current, err = migrations.Current(currentVersion); err != nil {
		return
	}
	if err = current.Down(db); err != nil {
		return
	}
	if err = current.Up(db); err != nil {
		return
	}

	return
}

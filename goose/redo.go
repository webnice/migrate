package goose

import (
	"database/sql"
)

// Redo Откат последней применённой миграции и её повторное выполнение.
func Redo(db *sql.DB, dir string) (err error) {
	var (
		currentVersion int64
		migrations     Migrations
		current        *Migration
	)

	if currentVersion, err = GetDBVersion(db); err != nil {
		return
	}
	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		return
	}
	if current, err = migrations.Current(currentVersion); err != nil {
		return
	}
	if err = current.Down(db); err != nil {
		return
	}
	err = current.Up(db)

	return
}

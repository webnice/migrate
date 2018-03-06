package goose

import (
	"database/sql"
	"log"
)

// Version prints the current version of the database.
func Version(db *sql.DB, dir string) (err error) {
	var current int64

	if current, err = GetDBVersion(db); err != nil {
		return
	}
	log.Printf("goose: version %014d\n", current)

	return
}

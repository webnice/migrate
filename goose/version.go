package goose

import (
	"database/sql"
	"log"
)

// Version Вывод текущей версии базы данных.
// Аргумент dir зарезервирован для совместимости интерфейса и не используется.
func Version(db *sql.DB, _ string) (err error) {
	const tplVersion = "goose: версия %014d"
	var current int64

	if current, err = GetDBVersion(db); err != nil {
		return
	}
	log.Printf(tplVersion, current)

	return
}

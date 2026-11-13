package goose

import (
	"database/sql"
	"log"
)

// Down Откат одной миграции от текущей версии.
func Down(db *sql.DB, dir string) (err error) {
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
		err = Errors().MigrationNotFound.Bind(currentVersion)
		return
	}
	err = current.Down(db)

	return
}

// DownTo Откат миграций до указанной версии.
func DownTo(db *sql.DB, dir string, version int64) (err error) {
	const tplNoMigrations = "goose: нет миграций для выполнения. Текущая версия: %d"
	var (
		migrations     Migrations
		currentVersion int64
		current        *Migration
		end            bool
	)

	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		return
	}
	for {
		if end {
			break
		}
		// Чтение текущей версии базы данных.
		if currentVersion, err = GetDBVersion(db); err != nil {
			end = true
			continue
		}
		// Поиск миграции текущей версии, её отсутствие означает завершение отката.
		if current, err = migrations.Current(currentVersion); err != nil {
			log.Printf(tplNoMigrations, currentVersion)
			err = nil
			end = true
			continue
		}
		if current.Version <= version {
			// Достигнута требуемая версия.
			log.Printf(tplNoMigrations, currentVersion)
			end = true
			continue
		}
		// Откат текущей миграции.
		if err = current.Down(db); err != nil {
			end = true
		}
	}

	return
}

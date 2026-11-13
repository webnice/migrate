package goose

import (
	"database/sql"
	"log"
)

// UpTo Применение миграций до указанной версии.
func UpTo(db *sql.DB, dir string, version int64) (err error) {
	const tplNoMigrations = "goose: нет миграций для выполнения. Текущая версия: %d"
	var (
		migrations Migrations
		current    int64
		next       *Migration
		end        bool
	)

	if migrations, err = CollectMigrations(dir, minVersion, version); err != nil {
		return
	}
	for {
		if end {
			break
		}
		// Чтение текущей версии базы данных.
		if current, err = GetDBVersion(db); err != nil {
			end = true
			continue
		}
		// Поиск следующей миграции.
		if next, err = migrations.Next(current); err != nil {
			if Errors().NoNextVersion.Is(err) {
				// Миграции закончились.
				log.Printf(tplNoMigrations, current)
				err = nil
			}
			end = true
			continue
		}
		// Применение следующей миграции.
		if err = next.Up(db); err != nil {
			end = true
		}
	}

	return
}

// Up Применение всех доступных миграций.
func Up(db *sql.DB, dir string) error { return UpTo(db, dir, MaxVersion) }

// UpByOne Применение одной следующей миграции.
func UpByOne(db *sql.DB, dir string) (err error) {
	const tplNoMigrations = "goose: нет миграций для выполнения. Текущая версия: %d"
	var (
		migrations     Migrations
		currentVersion int64
		next           *Migration
	)

	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		return
	}
	if currentVersion, err = GetDBVersion(db); err != nil {
		return
	}
	if next, err = migrations.Next(currentVersion); err != nil {
		if Errors().NoNextVersion.Is(err) {
			// Миграции закончились.
			log.Printf(tplNoMigrations, currentVersion)
			err = nil
		}
		return
	}
	err = next.Up(db)

	return
}
